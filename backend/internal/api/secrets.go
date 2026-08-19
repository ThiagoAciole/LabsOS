package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"labsos/backend/internal/safety"
)

var secretNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]{0,63}$`)

func secretsDir() string {
	if value := strings.TrimSpace(os.Getenv("LABSOS_SECRETS_DIR")); value != "" {
		return value
	}
	return "/var/lib/labsos/secrets"
}

func (s *server) secrets(w http.ResponseWriter, _ *http.Request) {
	entries, err := os.ReadDir(secretsDir())
	if os.IsNotExist(err) {
		respond(w, []map[string]any{}, nil, http.StatusOK)
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "SECRETS_LIST_FAILED", err.Error())
		return
	}
	result := make([]map[string]any, 0)
	for _, entry := range entries {
		if !entry.Type().IsRegular() || !secretNamePattern.MatchString(entry.Name()) {
			continue
		}
		info, infoErr := entry.Info()
		if infoErr != nil {
			continue
		}
		result = append(result, map[string]any{"name": entry.Name(), "updatedAt": info.ModTime()})
	}
	respond(w, result, nil, http.StatusOK)
}

func (s *server) putSecret(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if !secretNamePattern.MatchString(name) {
		writeError(w, http.StatusBadRequest, "INVALID_SECRET_NAME", "invalid secret name")
		return
	}
	var input struct {
		Value string `json:"value"`
	}
	if json.NewDecoder(r.Body).Decode(&input) != nil || input.Value == "" {
		writeError(w, http.StatusBadRequest, "INVALID_SECRET", "secret value is required")
		return
	}
	if !safety.RealOperationsEnabled() {
		s.auditStore.add("secret.put", name, "planned", "real operations disabled")
		respond(w, map[string]any{"name": name, "stored": false, "simulated": true, "message": "secret write planned; real operations are disabled"}, nil, http.StatusAccepted)
		return
	}
	if err := os.MkdirAll(secretsDir(), 0700); err != nil {
		writeError(w, http.StatusInternalServerError, "SECRETS_DIRECTORY_FAILED", err.Error())
		return
	}
	path := filepath.Join(secretsDir(), name)
	if err := os.WriteFile(path, []byte(input.Value+"\n"), 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "SECRET_WRITE_FAILED", err.Error())
		return
	}
	s.auditStore.add("secret.put", name, "success", "value omitted")
	respond(w, map[string]any{"name": name, "stored": true, "simulated": false}, nil, http.StatusOK)
}

func (s *server) deleteSecret(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if !secretNamePattern.MatchString(name) {
		writeError(w, http.StatusBadRequest, "INVALID_SECRET_NAME", "invalid secret name")
		return
	}
	if !safety.RealOperationsEnabled() {
		s.auditStore.add("secret.delete", name, "planned", "real operations disabled")
		respond(w, map[string]any{"name": name, "deleted": false, "simulated": true}, nil, http.StatusAccepted)
		return
	}
	if err := os.Remove(filepath.Join(secretsDir(), name)); err != nil && !os.IsNotExist(err) {
		writeError(w, http.StatusInternalServerError, "SECRET_DELETE_FAILED", err.Error())
		return
	}
	s.auditStore.add("secret.delete", name, "success", "")
	respond(w, map[string]any{"name": name, "deleted": true, "simulated": false}, nil, http.StatusOK)
}
