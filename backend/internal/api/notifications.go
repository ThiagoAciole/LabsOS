package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type notification struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
	Severity  string    `json:"severity"`
	CreatedAt time.Time `json:"createdAt"`
	Read      bool      `json:"read"`
}

type notificationStore struct {
	mu    sync.Mutex
	items []notification
	path  string
	next  uint64
}

func newNotificationStore() *notificationStore {
	path := os.Getenv("LABSOS_NOTIFICATIONS_FILE")
	if path == "" {
		path = "/tmp/labsos-notifications.json"
	}
	s := &notificationStore{path: path}
	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, &s.items)
	}
	for _, item := range s.items {
		if item.ID > "" {
			s.next++
		}
	}
	return s
}

func (s *notificationStore) add(title, message, source, severity string) notification {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	item := notification{ID: "notification-" + formatJobID(s.next), Type: severity, Title: title, Message: message, Source: source, Severity: severity, CreatedAt: time.Now().UTC()}
	s.items = append([]notification{item}, s.items...)
	if len(s.items) > 200 {
		s.items = s.items[:200]
	}
	s.persistLocked()
	return item
}

func (s *notificationStore) addOnce(title, message, source, severity string) (notification, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range s.items {
		if !item.Read && item.Title == title && item.Source == source {
			return item, false
		}
	}
	s.next++
	item := notification{ID: "notification-" + formatJobID(s.next), Type: severity, Title: title, Message: message, Source: source, Severity: severity, CreatedAt: time.Now().UTC()}
	s.items = append([]notification{item}, s.items...)
	if len(s.items) > 200 {
		s.items = s.items[:200]
	}
	s.persistLocked()
	return item, true
}
func (s *notificationStore) list() []notification {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]notification(nil), s.items...)
}
func (s *notificationStore) markRead(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.items {
		if s.items[i].ID == id {
			s.items[i].Read = true
			s.persistLocked()
			return true
		}
	}
	return false
}
func (s *notificationStore) delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.items {
		if s.items[i].ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			s.persistLocked()
			return true
		}
	}
	return false
}
func (s *notificationStore) persistLocked() {
	data, _ := json.MarshalIndent(s.items, "", "  ")
	_ = os.MkdirAll(filepath.Dir(s.path), 0750)
	_ = os.WriteFile(s.path, data, 0600)
}

func (s *server) notifications(w http.ResponseWriter, _ *http.Request) {
	respond(w, s.notificationStore.list(), nil, http.StatusOK)
}
func (s *server) readNotification(w http.ResponseWriter, r *http.Request) {
	if !s.notificationStore.markRead(r.PathValue("id")) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "notification not found")
		return
	}
	respond(w, map[string]any{"read": true}, nil, http.StatusOK)
}
func (s *server) deleteNotification(w http.ResponseWriter, r *http.Request) {
	if !s.notificationStore.delete(r.PathValue("id")) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "notification not found")
		return
	}
	s.auditStore.add("notification.delete", r.PathValue("id"), "success", "")
	respond(w, map[string]any{"deleted": true}, nil, http.StatusOK)
}
