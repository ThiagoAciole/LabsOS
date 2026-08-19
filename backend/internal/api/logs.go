package api

import (
	"context"
	"net/http"
	"os/exec"
	"strconv"
)

func (s *server) systemLogs(w http.ResponseWriter, r *http.Request) {
	respond(w, map[string]any{"source": "system", "logs": readJournal(r.Context(), "-n", journalLines(r))}, nil, http.StatusOK)
}
func (s *server) kernelLogs(w http.ResponseWriter, r *http.Request) {
	respond(w, map[string]any{"source": "kernel", "logs": readJournal(r.Context(), "-k", "-n", journalLines(r))}, nil, http.StatusOK)
}
func readJournal(ctx context.Context, args ...string) string {
	output, err := exec.CommandContext(ctx, "journalctl", append([]string{"--no-pager", "-o", "short-iso"}, args...)...).CombinedOutput()
	if err != nil {
		return "logs unavailable: " + string(output)
	}
	return string(output)
}
func journalLines(r *http.Request) string {
	value, err := strconv.Atoi(r.URL.Query().Get("lines"))
	if err != nil || value < 1 {
		return "100"
	}
	if value > 1000 {
		return "1000"
	}
	return strconv.Itoa(value)
}
