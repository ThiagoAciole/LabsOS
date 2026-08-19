package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"labsos/backend/internal/platform"
)

type declarativeApp struct {
	platform.App
	Compose   string    `json:"compose"`
	CreatedAt time.Time `json:"createdAt"`
}

type declarativeStore struct {
	mu   sync.Mutex
	path string
	apps []declarativeApp
}

type declarativeImportRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
	Version     string `json:"version"`
	Compose     string `json:"compose"`
}

var declarativeID = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,62}$`)

func newDeclarativeStore() *declarativeStore {
	path := os.Getenv("LABSOS_DECLARATIVE_APPS_FILE")
	if path == "" {
		path = "/tmp/labsos-declarative-apps.json"
	}
	store := &declarativeStore{path: path}
	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, &store.apps)
	}
	return store
}

func (s *declarativeStore) list() []declarativeApp {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]declarativeApp(nil), s.apps...)
}

func (s *declarativeStore) add(input declarativeImportRequest) (declarativeApp, error) {
	input.ID = strings.ToLower(strings.TrimSpace(input.ID))
	input.Name = strings.TrimSpace(input.Name)
	input.Compose = strings.TrimSpace(input.Compose)
	if !declarativeID.MatchString(input.ID) {
		return declarativeApp{}, fmt.Errorf("id must contain only lowercase letters, numbers, dots, underscores or hyphens")
	}
	if input.Name == "" || input.Compose == "" {
		return declarativeApp{}, fmt.Errorf("name and compose are required")
	}
	if len(input.Compose) > 512*1024 {
		return declarativeApp{}, fmt.Errorf("compose manifest is too large")
	}
	if !strings.Contains(input.Compose, "services:") && !strings.Contains(input.Compose, "services\n") {
		return declarativeApp{}, fmt.Errorf("compose manifest must declare services")
	}
	item := declarativeApp{App: platform.App{ID: input.ID, Kind: platform.AppKindUser, Name: input.Name, Description: input.Description, Icon: input.Icon, Category: input.Category, Source: "local", Version: input.Version, Installable: true, Actions: []string{"install"}}, Compose: input.Compose, CreatedAt: time.Now().UTC()}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.apps {
		if s.apps[i].ID == item.ID {
			s.apps[i] = item
			return item, s.persistLocked()
		}
	}
	s.apps = append(s.apps, item)
	return item, s.persistLocked()
}

func (s *declarativeStore) persistLocked() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0750); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.apps, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}

func (s *server) importDeclarativeApp(w http.ResponseWriter, r *http.Request) {
	var input declarativeImportRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 512*1024+4096)).Decode(&input); err != nil {
		writeError(w, 400, "INVALID_JSON", "declarative app must be valid JSON")
		return
	}
	item, err := s.declarative.add(input)
	if err != nil {
		writeError(w, 400, "INVALID_MANIFEST", err.Error())
		return
	}
	s.auditStore.add("app.import", item.ID, "success", "declarative manifest registered")
	s.eventHub.publish("catalog", "declarative app "+item.ID+" registered")
	respond(w, map[string]any{"app": item, "installed": false, "message": "manifest registered; installation requires a separate protected action"}, nil, http.StatusAccepted)
}
