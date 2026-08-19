package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"labsos/backend/internal/platform"
	"labsos/backend/internal/safety"
)

type serviceExposure struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Ports    []int  `json:"ports,omitempty"`
	URL      string `json:"url,omitempty"`
	Exposed  bool   `json:"exposed"`
	Provider string `json:"provider"`
	Conflict bool   `json:"conflict"`
}

type exposureStore struct {
	mu    sync.Mutex
	path  string
	items map[string]bool
}

func newExposureStore() *exposureStore {
	path := os.Getenv("LABSOS_EXPOSURES_FILE")
	if path == "" {
		path = "/tmp/labsos-exposures.json"
	}
	store := &exposureStore{path: path, items: map[string]bool{}}
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &store.items)
	}
	return store
}

func (s *exposureStore) get(id string) (bool, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.items[id]
	return value, ok
}

func (s *exposureStore) set(id string, exposed bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[id] = exposed
	data, err := json.MarshalIndent(s.items, "", "  ")
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0750); err != nil {
		return
	}
	tmp, err := os.CreateTemp(filepath.Dir(s.path), ".exposures-*.tmp")
	if err != nil {
		return
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err = tmp.Write(data); err == nil {
		err = tmp.Chmod(0600)
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err == nil {
		_ = os.Rename(tmpPath, s.path)
	}
}

func (s *server) services(w http.ResponseWriter, r *http.Request) {
	apps, err := s.provider.Apps(r.Context(), false)
	if err != nil {
		respond(w, nil, err, http.StatusOK)
		return
	}
	result := make([]serviceExposure, 0)
	owners := map[int]string{}
	conflicts := map[string]bool{}
	for _, app := range apps {
		for _, port := range app.Ports {
			if owner, exists := owners[port]; exists && owner != app.ID {
				conflicts[owner] = true
				conflicts[app.ID] = true
			} else {
				owners[port] = app.ID
			}
		}
	}
	for _, app := range apps {
		if len(app.Ports) == 0 && app.URL == "" {
			continue
		}
		exposed, known := s.exposures.get(app.ID)
		if !known {
			exposed = app.URL != ""
		}
		result = append(result, serviceExposure{ID: app.ID, Name: app.Name, Ports: app.Ports, URL: app.URL, Exposed: exposed, Provider: app.Source, Conflict: conflicts[app.ID]})
	}
	respond(w, result, nil, http.StatusOK)
}

func (s *server) updateServiceExposure(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "INVALID_SERVICE", "service id is required")
		return
	}
	var input struct {
		Exposed bool `json:"exposed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "exposure request must be valid JSON")
		return
	}
	if !safety.RealOperationsEnabled() {
		details := "service exposure is planned; no firewall or proxy was changed"
		if input.Exposed {
			details = "service exposure requested; no firewall or proxy was changed"
		}
		s.auditStore.add("service.exposure", id, "planned", details)
		s.eventHub.publish("service", id+" exposure planned: "+strconv.FormatBool(input.Exposed))
		respond(w, map[string]any{"id": id, "exposed": input.Exposed, "applied": false, "planned": true, "message": details}, nil, http.StatusAccepted)
		return
	}
	apps, err := s.provider.Apps(r.Context(), false)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "SERVICE_LOOKUP_FAILED", err.Error())
		return
	}
	var ports []int
	for _, app := range apps {
		if app.ID == id {
			ports = app.Ports
			break
		}
	}
	if len(ports) == 0 {
		writeError(w, http.StatusBadRequest, "SERVICE_PORTS_UNKNOWN", "service has no published TCP ports")
		return
	}
	if input.Exposed {
		if conflicts := servicePortConflicts(apps); conflicts[id] {
			writeError(w, http.StatusConflict, "PORT_CONFLICT", "service has a published port conflict")
			return
		}
	}
	if err := applyFirewallPorts(r.Context(), ports, input.Exposed); err != nil {
		s.auditStore.add("service.exposure", id, "error", err.Error())
		writeError(w, http.StatusBadGateway, "FIREWALL_UPDATE_FAILED", err.Error())
		return
	}
	details := "service exposure applied by firewalld"
	s.auditStore.add("service.exposure", id, "success", details)
	s.exposures.set(id, input.Exposed)
	s.eventHub.publish("service", id+" exposure changed: "+strconv.FormatBool(input.Exposed))
	respond(w, map[string]any{"id": id, "exposed": input.Exposed, "applied": true, "planned": false, "message": details}, nil, http.StatusAccepted)
}

func servicePortConflicts(apps []platform.App) map[string]bool {
	owners := map[int]string{}
	conflicts := map[string]bool{}
	for _, app := range apps {
		for _, port := range app.Ports {
			if owner, exists := owners[port]; exists && owner != app.ID {
				conflicts[owner] = true
				conflicts[app.ID] = true
			} else {
				owners[port] = app.ID
			}
		}
	}
	return conflicts
}

func applyFirewallPorts(ctx context.Context, ports []int, exposed bool) error {
	if _, err := exec.LookPath("firewall-cmd"); err == nil {
		if firewallErr := applyFirewalldPorts(ctx, ports, exposed); firewallErr == nil {
			return nil
		}
	}
	if _, err := exec.LookPath("nft"); err == nil {
		return applyNftablesPorts(ctx, ports, exposed)
	}
	return fmt.Errorf("no supported firewall backend found (firewall-cmd or nft)")
}

func applyFirewalldPorts(ctx context.Context, ports []int, exposed bool) error {
	for _, port := range ports {
		if port < 1 || port > 65535 {
			continue
		}
		action := "--remove-port=" + strconv.Itoa(port) + "/tcp"
		if exposed {
			action = "--add-port=" + strconv.Itoa(port) + "/tcp"
		}
		if output, err := exec.CommandContext(ctx, "firewall-cmd", "--permanent", action).CombinedOutput(); err != nil {
			return fmt.Errorf("firewall-cmd %s: %s", action, strings.TrimSpace(string(output)))
		}
	}
	if output, err := exec.CommandContext(ctx, "firewall-cmd", "--reload").CombinedOutput(); err != nil {
		return fmt.Errorf("firewall-cmd --reload: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

func applyNftablesPorts(ctx context.Context, ports []int, exposed bool) error {
	if len(ports) == 0 {
		return fmt.Errorf("no ports were provided")
	}
	if err := nftEnsure(ctx, "add", "table", "inet", "labsos"); err != nil {
		return fmt.Errorf("nftables table: %w", err)
	}
	if err := nftEnsure(ctx, "add", "set", "inet", "labsos", "exposed_tcp_ports", "{", "type", "inet_service", ";", "}"); err != nil {
		return fmt.Errorf("nftables set: %w", err)
	}
	chain, err := exec.CommandContext(ctx, "nft", "list", "chain", "inet", "filter", "input").CombinedOutput()
	if err != nil {
		return fmt.Errorf("nftables input chain unavailable: %s", strings.TrimSpace(string(chain)))
	}
	if !strings.Contains(string(chain), "@exposed_tcp_ports") {
		if output, ruleErr := exec.CommandContext(ctx, "nft", "insert", "rule", "inet", "filter", "input", "tcp", "dport", "@exposed_tcp_ports", "accept").CombinedOutput(); ruleErr != nil {
			return fmt.Errorf("nftables input rule: %s", strings.TrimSpace(string(output)))
		}
	}
	for _, port := range ports {
		if port < 1 || port > 65535 {
			continue
		}
		element := "{" + strconv.Itoa(port) + "}"
		action := "add"
		if !exposed {
			action = "delete"
		}
		if output, commandErr := exec.CommandContext(ctx, "nft", action, "element", "inet", "labsos", "exposed_tcp_ports", element).CombinedOutput(); commandErr != nil && exposed {
			return fmt.Errorf("nftables %s port %d: %s", action, port, strings.TrimSpace(string(output)))
		}
	}
	return nil
}

func nftEnsure(ctx context.Context, args ...string) error {
	output, err := exec.CommandContext(ctx, "nft", args...).CombinedOutput()
	if err != nil && !strings.Contains(strings.ToLower(string(output)), "exists") {
		return fmt.Errorf("%s: %s", strings.Join(args, " "), strings.TrimSpace(string(output)))
	}
	return nil
}
