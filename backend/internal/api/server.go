package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"labsos/backend/internal/platform"
)

type server struct {
	provider platform.Provider
	mux      *http.ServeMux
}

func New(provider platform.Provider) http.Handler {
	s := &server{provider: provider, mux: http.NewServeMux()}
	s.mux.HandleFunc("GET /api/v1/system/summary", s.systemSummary)
	s.mux.HandleFunc("GET /api/v1/system/health", s.systemHealth)
	s.mux.HandleFunc("POST /api/v1/system/{action}", s.power)
	s.mux.HandleFunc("GET /api/v1/system/update", s.updateStatus)
	s.mux.HandleFunc("POST /api/v1/system/update", s.applyUpdate)
	s.mux.HandleFunc("GET /api/v1/apps", s.apps)
	s.mux.HandleFunc("GET /api/v1/catalog/apps", s.catalog)
	s.mux.HandleFunc("GET /api/v1/catalog/sources", s.catalogSources)
	s.mux.HandleFunc("POST /api/v1/apps/{id}/{action}", s.appAction)
	s.mux.HandleFunc("DELETE /api/v1/apps/{id}", s.removeApp)
	s.mux.HandleFunc("GET /api/v1/apps/{id}/logs", s.appLogs)
	s.mux.HandleFunc("GET /api/v1/settings", s.settings)
	s.mux.HandleFunc("PUT /api/v1/settings/{section}", s.updateSettings)
	s.mux.HandleFunc("GET /api/v1/jobs/{id}", s.job)
	s.mux.HandleFunc("GET /api/v1/events", s.events)
	s.mux.HandleFunc("/api/v1/", func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "endpoint not found")
	})
	return s.mux
}

func (s *server) systemSummary(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.SystemSummary(r.Context())
	respond(w, value, err, http.StatusOK)
}

func (s *server) systemHealth(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.SystemHealth(r.Context())
	respond(w, value, err, http.StatusOK)
}

func (s *server) power(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.Power(r.Context(), r.PathValue("action"))
	respond(w, value, err, http.StatusAccepted)
}

func (s *server) updateStatus(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.UpdateStatus(r.Context())
	respond(w, value, err, http.StatusOK)
}
func (s *server) applyUpdate(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.ApplyUpdate(r.Context())
	respond(w, value, err, http.StatusAccepted)
}

func (s *server) apps(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.Apps(r.Context(), false)
	respond(w, value, err, http.StatusOK)
}

func (s *server) catalog(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.Apps(r.Context(), true)
	respond(w, value, err, http.StatusOK)
}

func (s *server) catalogSources(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.CatalogSources(r.Context())
	respond(w, value, err, http.StatusOK)
}

func (s *server) appAction(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.AppAction(r.Context(), r.PathValue("id"), r.PathValue("action"))
	respond(w, value, err, http.StatusAccepted)
}

func (s *server) removeApp(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.RemoveApp(r.Context(), r.PathValue("id"))
	respond(w, value, err, http.StatusAccepted)
}

func (s *server) appLogs(w http.ResponseWriter, r *http.Request) {
	lines := 100
	if value := r.URL.Query().Get("lines"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			lines = parsed
		}
	}
	if lines < 1 {
		lines = 1
	}
	if lines > 1000 {
		lines = 1000
	}
	value, err := s.provider.AppLogs(r.Context(), r.PathValue("id"), lines)
	respond(w, map[string]string{"app": r.PathValue("id"), "logs": value}, err, http.StatusOK)
}

func (s *server) settings(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.Settings(r.Context())
	respond(w, value, err, http.StatusOK)
}

func (s *server) updateSettings(w http.ResponseWriter, r *http.Request) {
	var patch map[string]any
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	if err := decoder.Decode(&patch); err != nil || patch == nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "request body must be a JSON object")
		return
	}
	value, err := s.provider.UpdateSettings(r.Context(), r.PathValue("section"), patch)
	respond(w, value, err, http.StatusOK)
}

func (s *server) events(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.Events(r.Context())
	respond(w, value, err, http.StatusOK)
}

func (s *server) job(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.Job(r.Context(), r.PathValue("id"))
	respond(w, value, err, http.StatusOK)
}

func respond(w http.ResponseWriter, value any, err error, status int) {
	if err != nil {
		writeProviderError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeProviderError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, platform.ErrNotFound):
		writeError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
	case errors.Is(err, platform.ErrUnsupported):
		writeError(w, http.StatusBadRequest, "UNSUPPORTED_OPERATION", "operation is not supported")
	case errors.Is(err, platform.ErrUnavailable):
		writeError(w, http.StatusServiceUnavailable, "PROVIDER_UNAVAILABLE", "active provider is unavailable")
	default:
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]any{"code": code, "message": message, "details": map[string]any{}}})
}
