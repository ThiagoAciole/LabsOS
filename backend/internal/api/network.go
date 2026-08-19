package api

import (
	"context"
	"encoding/json"
	"labsos/backend/internal/safety"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type networkSnapshot struct {
	Interfaces any      `json:"interfaces"`
	Routes     string   `json:"routes"`
	DNS        []string `json:"dns"`
	Simulated  bool     `json:"simulated"`
}

type wifiSnapshot struct {
	Available bool     `json:"available"`
	Devices   []string `json:"devices"`
	Networks  []string `json:"networks"`
	Current   string   `json:"current,omitempty"`
	Simulated bool     `json:"simulated"`
}

func (s *server) network(w http.ResponseWriter, r *http.Request) {
	interfaces := []any{}
	if output, err := exec.CommandContext(r.Context(), "ip", "-j", "address").Output(); err == nil {
		_ = json.Unmarshal(output, &interfaces)
	}
	dns := []string{}
	if data, err := os.ReadFile("/etc/resolv.conf"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			fields := strings.Fields(line)
			if len(fields) >= 2 && fields[0] == "nameserver" {
				dns = append(dns, fields[1])
			}
		}
	}
	respond(w, networkSnapshot{Interfaces: interfaces, Routes: commandText(r.Context(), "ip", "route"), DNS: dns, Simulated: !safety.NetworkChangesEnabled()}, nil, http.StatusOK)
}

func (s *server) updateNetwork(w http.ResponseWriter, r *http.Request) {
	var patch map[string]any
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "network patch must be JSON")
		return
	}
	interfaceName, _ := patch["interface"].(string)
	dhcp, hasDHCP := patch["dhcp"].(bool)
	if !safety.NetworkChangesEnabled() || interfaceName == "" || !hasDHCP {
		s.auditStore.add("network.update", interfaceName, "planned", "network changes disabled or unsupported request")
		respond(w, map[string]any{"applied": false, "simulated": true, "message": "network changes are disabled by safety policy or request is unsupported", "requested": patch}, nil, http.StatusAccepted)
		return
	}
	if !validInterfaceName(interfaceName) {
		writeError(w, http.StatusBadRequest, "INVALID_INTERFACE", "invalid network interface")
		return
	}
	method := "manual"
	if dhcp {
		method = "auto"
	}
	args := []string{"connection", "modify", interfaceName, "ipv4.method", method}
	if !dhcp {
		address, _ := patch["address"].(string)
		gateway, _ := patch["gateway"].(string)
		dns, _ := patch["dns"].(string)
		if _, _, err := net.ParseCIDR(address); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_ADDRESS", "static address must be CIDR")
			return
		}
		if net.ParseIP(gateway) == nil {
			writeError(w, http.StatusBadRequest, "INVALID_GATEWAY", "gateway must be an IP address")
			return
		}
		for _, item := range strings.Split(dns, ",") {
			if net.ParseIP(strings.TrimSpace(item)) == nil {
				writeError(w, http.StatusBadRequest, "INVALID_DNS", "DNS entries must be IP addresses")
				return
			}
		}
		args = append(args, "ipv4.addresses", address, "ipv4.gateway", gateway, "ipv4.dns", dns)
	}
	if output, err := exec.CommandContext(r.Context(), "nmcli", args...).CombinedOutput(); err != nil {
		details := strings.TrimSpace(string(output))
		s.auditStore.add("network.update", interfaceName, "error", details)
		writeError(w, http.StatusBadGateway, "NETWORK_UPDATE_FAILED", details)
		return
	}
	if output, err := exec.CommandContext(r.Context(), "nmcli", "connection", "up", interfaceName).CombinedOutput(); err != nil {
		details := strings.TrimSpace(string(output))
		s.auditStore.add("network.update", interfaceName, "error", details)
		writeError(w, http.StatusBadGateway, "NETWORK_ACTIVATION_FAILED", details)
		return
	}
	s.auditStore.add("network.update", interfaceName, "success", "dhcp="+strconv.FormatBool(dhcp))
	respond(w, map[string]any{"applied": true, "simulated": false, "message": "network connection updated", "interface": interfaceName, "dhcp": dhcp}, nil, http.StatusAccepted)
}

func validInterfaceName(value string) bool {
	if value == "" || len(value) > 15 {
		return false
	}
	for _, char := range value {
		if !(char == '-' || char == '_' || char == '.' || char >= '0' && char <= '9' || char >= 'A' && char <= 'Z' || char >= 'a' && char <= 'z') {
			return false
		}
	}
	return true
}

func (s *server) wifi(w http.ResponseWriter, r *http.Request) {
	result := wifiSnapshot{Devices: []string{}, Networks: []string{}, Simulated: true}
	if output, err := exec.CommandContext(r.Context(), "nmcli", "-t", "-f", "DEVICE,TYPE,STATE,CONNECTION", "device").Output(); err == nil {
		result.Available = true
		for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
			fields := strings.Split(line, ":")
			if len(fields) < 4 || fields[1] != "wifi" {
				continue
			}
			result.Devices = append(result.Devices, fields[0])
			if fields[2] == "connected" && fields[3] != "--" {
				result.Current = fields[3]
			}
		}
		if output, err := exec.CommandContext(r.Context(), "nmcli", "-t", "-f", "SSID,SECURITY,SIGNAL", "device", "wifi", "list", "--rescan", "no").Output(); err == nil {
			for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
				if fields := strings.Split(line, ":"); len(fields) > 0 && fields[0] != "" {
					result.Networks = append(result.Networks, fields[0])
				}
			}
		}
	}
	respond(w, result, nil, http.StatusOK)
}

func (s *server) updateWiFi(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Action   string `json:"action"`
		SSID     string `json:"ssid"`
		Password string `json:"password"`
		Device   string `json:"device"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "wifi request must be JSON")
		return
	}
	input.SSID = strings.TrimSpace(input.SSID)
	if input.Action != "connect" && input.Action != "reconnect" && input.Action != "forget" {
		writeError(w, http.StatusBadRequest, "INVALID_WIFI_ACTION", "unsupported Wi-Fi action")
		return
	}
	if input.SSID == "" || len(input.SSID) > SSIDLimit || strings.ContainsAny(input.SSID, "\r\n") {
		writeError(w, http.StatusBadRequest, "INVALID_SSID", "invalid SSID")
		return
	}
	if !safety.NetworkChangesEnabled() {
		message := "Wi-Fi change planned; network changes are disabled"
		s.auditStore.add("wifi."+input.Action, input.SSID, "planned", message)
		respond(w, map[string]any{"applied": false, "simulated": true, "message": message}, nil, http.StatusAccepted)
		return
	}
	var cmd *exec.Cmd
	switch input.Action {
	case "connect":
		if input.Password == "" {
			writeError(w, http.StatusBadRequest, "PASSWORD_REQUIRED", "Wi-Fi password is required")
			return
		}
		args := []string{"device", "wifi", "connect", input.SSID, "password", input.Password}
		if input.Device != "" {
			if !validInterfaceName(input.Device) {
				writeError(w, http.StatusBadRequest, "INVALID_INTERFACE", "invalid Wi-Fi device")
				return
			}
			args = append(args, "ifname", input.Device)
		}
		cmd = exec.CommandContext(r.Context(), "nmcli", args...)
	case "reconnect":
		cmd = exec.CommandContext(r.Context(), "nmcli", "connection", "up", input.SSID)
	case "forget":
		cmd = exec.CommandContext(r.Context(), "nmcli", "connection", "delete", input.SSID)
	}
	if output, err := cmd.CombinedOutput(); err != nil {
		details := strings.TrimSpace(string(output))
		s.auditStore.add("wifi."+input.Action, input.SSID, "error", details)
		writeError(w, http.StatusBadGateway, "WIFI_UPDATE_FAILED", details)
		return
	}
	s.auditStore.add("wifi."+input.Action, input.SSID, "success", "")
	respond(w, map[string]any{"applied": true, "simulated": false, "message": "Wi-Fi updated"}, nil, http.StatusAccepted)
}

const SSIDLimit = 128

func commandText(ctx context.Context, command string, args ...string) string {
	output, err := exec.CommandContext(ctx, command, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
