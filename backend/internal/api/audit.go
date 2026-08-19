package api

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"
)

type auditEntry struct {
	ID      string    `json:"id"`
	At      time.Time `json:"at"`
	Actor   string    `json:"actor"`
	Action  string    `json:"action"`
	Target  string    `json:"target,omitempty"`
	Status  string    `json:"status"`
	Details string    `json:"details,omitempty"`
}
type auditStore struct {
	mu      sync.Mutex
	entries []auditEntry
	seq     int
}

func newAuditStore() *auditStore {
	s := &auditStore{}
	if data, err := os.ReadFile(auditPath()); err == nil {
		_ = json.Unmarshal(data, &s.entries)
		s.seq = len(s.entries)
	}
	return s
}
func auditPath() string {
	if value := os.Getenv("LABSOS_AUDIT_FILE"); value != "" {
		return value
	}
	return "/tmp/labsos-audit.json"
}
func (s *auditStore) add(action, target, status, details string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	s.entries = append(s.entries, auditEntry{ID: "audit-" + formatJobID(uint64(s.seq)), At: time.Now().UTC(), Actor: "local-session", Action: action, Target: target, Status: status, Details: details})
	if len(s.entries) > 1000 {
		s.entries = s.entries[len(s.entries)-1000:]
	}
	data, _ := json.MarshalIndent(s.entries, "", "  ")
	_ = os.WriteFile(auditPath(), data, 0600)
}
func (s *auditStore) list() []auditEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]auditEntry, len(s.entries))
	copy(result, s.entries)
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}
func (s *server) audit(w http.ResponseWriter, _ *http.Request) {
	respond(w, s.auditStore.list(), nil, http.StatusOK)
}
