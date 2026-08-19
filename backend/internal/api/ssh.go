package api

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"labsos/backend/internal/safety"
)

type sshKey struct {
	Fingerprint string `json:"fingerprint"`
	Comment     string `json:"comment,omitempty"`
}
type sshStatus struct {
	Enabled   bool     `json:"enabled"`
	Active    bool     `json:"active"`
	Port      int      `json:"port"`
	Keys      []sshKey `json:"keys"`
	Simulated bool     `json:"simulated"`
}

func (s *server) ssh(w http.ResponseWriter, r *http.Request) {
	active := exec.CommandContext(r.Context(), "systemctl", "is-active", "--quiet", "ssh").Run() == nil
	port := 22
	if output, err := exec.CommandContext(r.Context(), "sshd", "-T").Output(); err == nil {
		for _, line := range strings.Split(string(output), "\n") {
			fields := strings.Fields(line)
			if len(fields) == 2 && fields[0] == "port" {
				if value, parseErr := strconv.Atoi(fields[1]); parseErr == nil {
					port = value
				}
			}
		}
	}
	keys := readSSHKeys()
	respond(w, sshStatus{Enabled: active, Active: active, Port: port, Keys: keys, Simulated: !safety.RealOperationsEnabled()}, nil, http.StatusOK)
}

func (s *server) updateSSH(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Enabled     *bool  `json:"enabled"`
		Action      string `json:"action"`
		Key         string `json:"key"`
		Fingerprint string `json:"fingerprint"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid SSH request")
		return
	}
	if input.Action == "add-key" || input.Action == "remove-key" {
		if input.Action == "add-key" && !validPublicKey(input.Key) {
			writeError(w, http.StatusBadRequest, "INVALID_KEY", "a valid public key is required")
			return
		}
		if input.Action == "remove-key" && strings.TrimSpace(input.Fingerprint) == "" {
			writeError(w, http.StatusBadRequest, "INVALID_KEY", "fingerprint is required")
			return
		}
		if !safety.RealOperationsEnabled() {
			s.auditStore.add("ssh."+input.Action, input.Fingerprint, "planned", "real operations disabled")
			respond(w, map[string]any{"applied": false, "simulated": true, "message": "SSH key change is planned; real operations are disabled"}, nil, http.StatusAccepted)
			return
		}
		if err := changeSSHKeys(input.Action, input.Key, input.Fingerprint); err != nil {
			writeError(w, http.StatusBadGateway, "SSH_KEY_UPDATE_FAILED", err.Error())
			return
		}
		s.auditStore.add("ssh."+input.Action, input.Fingerprint, "success", "")
		respond(w, map[string]any{"applied": true, "simulated": false, "message": "SSH keys updated"}, nil, http.StatusAccepted)
		return
	}
	if input.Enabled == nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "enabled or key action is required")
		return
	}
	if !safety.RealOperationsEnabled() {
		s.auditStore.add("ssh.update", "ssh", "planned", "real operations disabled")
		respond(w, map[string]any{"applied": false, "simulated": true, "message": "SSH changes are disabled by safety policy", "requested": input.Enabled}, nil, http.StatusAccepted)
		return
	}
	action := "disable"
	if *input.Enabled {
		action = "enable"
	}
	cmd := exec.CommandContext(r.Context(), "systemctl", action, "--now", "ssh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		details := strings.TrimSpace(string(output))
		s.auditStore.add("ssh.update", "ssh", "error", details)
		writeError(w, http.StatusBadGateway, "SSH_UPDATE_FAILED", details)
		return
	}
	s.auditStore.add("ssh.update", "ssh", "success", action)
	respond(w, map[string]any{"applied": true, "simulated": false, "message": "SSH " + action + "d"}, nil, http.StatusAccepted)
}

func validPublicKey(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" || len(value) > 4096 || strings.ContainsAny(value, "\r\n") {
		return false
	}
	fields := strings.Fields(value)
	return len(fields) >= 2 && (strings.HasPrefix(fields[0], "ssh-") || fields[0] == "ecdsa-sha2-nistp256" || fields[0] == "sk-ssh-ed25519@openssh.com")
}

func changeSSHKeys(action, key, fingerprint string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, ".ssh")
	path := filepath.Join(dir, "authorized_keys")
	data, _ := os.ReadFile(path)
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if action == "add-key" {
		lines = append(lines, strings.TrimSpace(key))
	} else {
		filtered := lines[:0]
		for _, line := range lines {
			if publicKeyFingerprint(strings.TrimSpace(line)) != fingerprint {
				filtered = append(filtered, line)
			}
		}
		lines = filtered
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(strings.TrimRight(strings.Join(lines, "\n"), "\n")+"\n"), 0600)
}
func readSSHKeys() []sshKey {
	home, err := os.UserHomeDir()
	if err != nil {
		return []sshKey{}
	}
	path := filepath.Join(home, ".ssh", "authorized_keys")
	file, err := os.Open(path)
	if err != nil {
		return []sshKey{}
	}
	defer file.Close()
	var result []sshKey
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) < 2 || strings.HasPrefix(fields[0], "#") {
			continue
		}
		fingerprint := publicKeyFingerprint(line)
		if fingerprint == "" {
			continue
		}
		comment := ""
		if len(fields) > 2 {
			comment = strings.Join(fields[2:], " ")
		}
		result = append(result, sshKey{Fingerprint: fingerprint, Comment: comment})
	}
	return result
}

func publicKeyFingerprint(line string) string {
	command := exec.Command("ssh-keygen", "-lf", "-")
	command.Stdin = strings.NewReader(line + "\n")
	output, err := command.Output()
	if err != nil {
		return ""
	}
	fields := strings.Fields(string(output))
	if len(fields) < 2 {
		return ""
	}
	return fields[1]
}
