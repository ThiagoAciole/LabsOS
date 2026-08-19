package api

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type diagnosticCheck struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (s *server) diagnostics(w http.ResponseWriter, r *http.Request) {
	checks := []diagnosticCheck{
		commandCheck(r.Context(), "docker", "docker", "info"),
		commandCheck(r.Context(), "labs-api", "systemctl", "is-active", "labs-api"),
		commandCheck(r.Context(), "labsd", "systemctl", "is-active", "labsd"),
	}
	if info, err := os.Stat("/DATA"); err == nil && info.IsDir() {
		checks = append(checks, diagnosticCheck{ID: "data", Status: "healthy", Message: "/DATA disponível"})
	} else {
		checks = append(checks, diagnosticCheck{ID: "data", Status: "error", Message: "/DATA indisponível"})
	}
	checks = append(checks, commandCheck(r.Context(), "internet", "curl", "-fsS", "--max-time", "3", "https://deb.debian.org/"))
	checks = append(checks, commandCheck(r.Context(), "dns", "getent", "hosts", "deb.debian.org"))
	checks = append(checks, commandCheck(r.Context(), "gateway", "sh", "-c", "ip route | grep -q '^default'"))
	respond(w, map[string]any{"checks": checks, "healthy": diagnosticsHealthy(checks)}, nil, http.StatusOK)
}

func commandCheck(ctx context.Context, id string, command string, args ...string) diagnosticCheck {
	output, err := exec.CommandContext(ctx, command, args...).CombinedOutput()
	if err != nil {
		return diagnosticCheck{ID: id, Status: "error", Message: strings.TrimSpace(string(output))}
	}
	return diagnosticCheck{ID: id, Status: "healthy", Message: strings.TrimSpace(string(output))}
}
func diagnosticsHealthy(checks []diagnosticCheck) bool {
	for _, check := range checks {
		if check.Status == "error" {
			return false
		}
	}
	return true
}
