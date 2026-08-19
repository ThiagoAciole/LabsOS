package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"labsos/backend/internal/platform"
	"labsos/backend/internal/safety"
)

type server struct {
	provider          platform.Provider
	mux               *http.ServeMux
	installerMu       sync.Mutex
	installerJob      *installerJob
	jobs              *jobEngine
	notificationStore *notificationStore
	sessions          *sessionStore
	scheduler         *schedulerStore
	auditStore        *auditStore
	declarative       *declarativeStore
	exposures         *exposureStore
	eventHub          *eventHub
}

func New(provider platform.Provider) http.Handler {
	eventHub := newEventHub()
	s := &server{provider: provider, mux: http.NewServeMux(), jobs: newJobEngine(), notificationStore: newNotificationStore(), sessions: newSessionStore(), scheduler: newSchedulerStore(), auditStore: newAuditStore(), declarative: newDeclarativeStore(), exposures: newExposureStore(), eventHub: eventHub}
	s.jobs.setPublisher(func(message string) { eventHub.publish("job", message) })
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			s.scheduler.tick(context.Background(), s)
			s.monitorHealth(context.Background())
		}
	}()
	s.mux.HandleFunc("GET /api/v1/system/summary", s.systemSummary)
	s.mux.HandleFunc("GET /api/v1/system/health", s.systemHealth)
	s.mux.HandleFunc("GET /api/v1/system/metrics", s.metrics)
	s.mux.HandleFunc("POST /api/v1/system/{action}", s.power)
	s.mux.HandleFunc("GET /api/v1/system/update", s.updateStatus)
	s.mux.HandleFunc("POST /api/v1/system/update", s.applyUpdate)
	s.mux.HandleFunc("POST /api/v1/system/rollback", s.rollbackUpdate)
	s.mux.HandleFunc("GET /api/v1/apps", s.apps)
	s.mux.HandleFunc("GET /api/v1/catalog/apps", s.catalog)
	s.mux.HandleFunc("GET /api/v1/catalog/sources", s.catalogSources)
	s.mux.HandleFunc("POST /api/v1/catalog/apps/import", s.importDeclarativeApp)
	s.mux.HandleFunc("POST /api/v1/apps/{id}/{action}", s.appAction)
	s.mux.HandleFunc("DELETE /api/v1/apps/{id}", s.removeApp)
	s.mux.HandleFunc("GET /api/v1/apps/{id}/logs", s.appLogs)
	s.mux.HandleFunc("GET /api/v1/apps/{id}/health", s.appHealth)
	s.mux.HandleFunc("GET /api/v1/apps/{id}/metrics", s.appMetrics)
	s.mux.HandleFunc("GET /api/v1/services", s.services)
	s.mux.HandleFunc("PUT /api/v1/services/{id}/exposure", s.updateServiceExposure)
	s.mux.HandleFunc("GET /api/v1/logs/system", s.systemLogs)
	s.mux.HandleFunc("GET /api/v1/logs/kernel", s.kernelLogs)
	s.mux.HandleFunc("GET /api/v1/settings", s.settings)
	s.mux.HandleFunc("PUT /api/v1/settings/{section}", s.updateSettings)
	s.mux.HandleFunc("GET /api/v1/network", s.network)
	s.mux.HandleFunc("GET /api/v1/network/wifi", s.wifi)
	s.mux.HandleFunc("PUT /api/v1/network/wifi", s.updateWiFi)
	s.mux.HandleFunc("GET /api/v1/storage", s.storage)
	s.mux.HandleFunc("GET /api/v1/files", s.files)
	s.mux.HandleFunc("PUT /api/v1/network", s.updateNetwork)
	s.mux.HandleFunc("GET /api/v1/diagnostics", s.diagnostics)
	s.mux.HandleFunc("GET /api/v1/access/ssh", s.ssh)
	s.mux.HandleFunc("PUT /api/v1/access/ssh", s.updateSSH)
	s.mux.HandleFunc("GET /api/v1/secrets", s.secrets)
	s.mux.HandleFunc("PUT /api/v1/secrets/{name}", s.putSecret)
	s.mux.HandleFunc("DELETE /api/v1/secrets/{name}", s.deleteSecret)
	s.mux.HandleFunc("GET /api/v1/backups", s.backups)
	s.mux.HandleFunc("GET /api/v1/backups/{id}/verify", s.verifyBackup)
	s.mux.HandleFunc("POST /api/v1/backups/{id}/restore", s.restoreBackup)
	s.mux.HandleFunc("DELETE /api/v1/backups/{id}", s.deleteBackup)
	s.mux.HandleFunc("POST /api/v1/backups", s.createBackup)
	s.mux.HandleFunc("GET /api/v1/auth/status", s.authStatus)
	s.mux.HandleFunc("POST /api/v1/auth/login", s.login)
	s.mux.HandleFunc("POST /api/v1/auth/logout", s.logout)
	s.mux.HandleFunc("POST /api/v1/auth/password", s.changePassword)
	s.mux.HandleFunc("GET /api/v1/auth/sessions", s.sessionsList)
	s.mux.HandleFunc("DELETE /api/v1/auth/sessions/{id}", s.revokeSession)
	s.mux.HandleFunc("GET /api/v1/jobs/{id}", s.job)
	s.mux.HandleFunc("GET /api/v1/jobs", s.jobsList)
	s.mux.HandleFunc("POST /api/v1/jobs/{id}/cancel", s.cancelJob)
	s.mux.HandleFunc("GET /api/v1/scheduler/tasks", s.schedulerList)
	s.mux.HandleFunc("POST /api/v1/scheduler/tasks", s.schedulerCreate)
	s.mux.HandleFunc("DELETE /api/v1/scheduler/tasks/{id}", s.schedulerDelete)
	s.mux.HandleFunc("GET /api/v1/audit", s.audit)
	s.mux.HandleFunc("GET /api/v1/events", s.events)
	s.mux.HandleFunc("GET /api/v1/events/stream", s.eventsStream)
	s.mux.HandleFunc("GET /api/v1/notifications", s.notifications)
	s.mux.HandleFunc("POST /api/v1/notifications/{id}/read", s.readNotification)
	s.mux.HandleFunc("DELETE /api/v1/notifications/{id}", s.deleteNotification)
	s.mux.HandleFunc("GET /api/v1/installer/status", s.installerStatus)
	s.mux.HandleFunc("GET /api/v1/installer/disks", s.installerDisks)
	s.mux.HandleFunc("POST /api/v1/installer/validate", s.installerValidate)
	s.mux.HandleFunc("POST /api/v1/installer/start", s.installerStart)
	s.mux.HandleFunc("GET /api/v1/installer/jobs/{id}", s.installerJobStatus)
	s.mux.HandleFunc("POST /api/v1/installer/jobs/{id}/cancel", s.installerCancel)
	s.mux.HandleFunc("GET /api/v1/installer/events", s.installerEvents)
	s.mux.HandleFunc("POST /api/v1/installer/reboot", s.installerReboot)
	s.mux.HandleFunc("/api/v1/", func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "endpoint not found")
	})
	return s.authMiddleware(s.mux)
}

func (s *server) monitorHealth(ctx context.Context) {
	summary, err := s.provider.SystemSummary(ctx)
	if err != nil {
		return
	}
	if summary.StorageTotal > 0 && float64(summary.StorageUsed)/float64(summary.StorageTotal) >= 0.9 {
		if _, added := s.notificationStore.addOnce("Pouco espaço", "O armazenamento do sistema está acima de 90%.", "system.storage", "warning"); added {
			s.eventHub.publish("notification", "Pouco espaço disponível")
		}
	}
	if summary.Temperature >= 85 {
		if _, added := s.notificationStore.addOnce("Temperatura alta", "A temperatura do sistema está acima de 85 °C.", "system.temperature", "error"); added {
			s.eventHub.publish("notification", "Temperatura alta")
		}
	}
}

func (s *server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		passwordSetup := r.URL.Path == "/api/v1/auth/password" && configuredPassword() == ""
		if r.URL.Path == "/api/v1/auth/status" || r.URL.Path == "/api/v1/auth/login" || r.URL.Path == "/api/v1/auth/logout" || passwordSetup || configuredPassword() == "" {
			next.ServeHTTP(w, r)
			return
		}
		cookie, err := r.Cookie("labsos_session")
		if err != nil || !s.sessions.valid(cookie.Value) {
			writeError(w, http.StatusUnauthorized, "AUTH_REQUIRED", "authentication required")
			return
		}
		next.ServeHTTP(w, r)
	})
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
	action := r.PathValue("action")
	if action != "reboot" && action != "shutdown" {
		writeError(w, http.StatusBadRequest, "INVALID_POWER_ACTION", "unsupported power action")
		return
	}
	if !safety.RealOperationsEnabled() {
		message := action + " planned; real operations are disabled"
		s.auditStore.add("system."+action, "labsos", "planned", message)
		s.eventHub.publish("system", message)
		respond(w, map[string]any{"id": "planned-" + action, "status": "success", "message": message, "simulated": true}, nil, http.StatusAccepted)
		return
	}
	value, err := s.provider.Power(r.Context(), action)
	if err != nil {
		s.auditStore.add("system."+action, "labsos", "error", err.Error())
	} else {
		s.auditStore.add("system."+action, "labsos", "success", "")
		s.eventHub.publish("system", action+" accepted")
	}
	respond(w, value, err, http.StatusAccepted)
}

func (s *server) updateStatus(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.UpdateStatus(r.Context())
	respond(w, value, err, http.StatusOK)
}
func (s *server) applyUpdate(w http.ResponseWriter, r *http.Request) {
	if !safety.RealOperationsEnabled() {
		s.auditStore.add("system.update", "labsos", "planned", "real operations disabled")
		respond(w, map[string]any{"currentVersion": "", "latestVersion": "", "updateAvailable": false, "applied": false, "simulated": true, "message": "update planned; real operations are disabled"}, nil, http.StatusAccepted)
		return
	}
	value, err := s.provider.ApplyUpdate(r.Context())
	if err != nil {
		s.auditStore.add("system.update", "labsos", "error", err.Error())
	} else {
		s.auditStore.add("system.update", "labsos", "success", "")
	}
	respond(w, value, err, http.StatusAccepted)
}

func (s *server) rollbackUpdate(w http.ResponseWriter, r *http.Request) {
	if !safety.RealOperationsEnabled() {
		s.auditStore.add("system.rollback", "labsos", "planned", "real operations disabled")
		respond(w, map[string]any{"applied": false, "simulated": true, "message": "rollback planned; real operations are disabled"}, nil, http.StatusAccepted)
		return
	}
	value, err := s.provider.RollbackUpdate(r.Context())
	if err != nil {
		s.auditStore.add("system.rollback", "labsos", "error", err.Error())
		respond(w, nil, err, http.StatusServiceUnavailable)
		return
	}
	s.auditStore.add("system.rollback", "labsos", "success", value.CurrentVersion)
	s.eventHub.publish("system", "rollback completed: "+value.CurrentVersion)
	respond(w, map[string]any{"applied": true, "simulated": false, "status": value}, nil, http.StatusAccepted)
}

func (s *server) apps(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.Apps(r.Context(), false)
	respond(w, value, err, http.StatusOK)
}

func (s *server) catalog(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.Apps(r.Context(), true)
	if err == nil {
		for _, item := range s.declarative.list() {
			value = append(value, item.App)
		}
		query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
		category := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("category")))
		if query != "" || category != "" {
			apps := value[:0]
			for _, app := range value {
				matchesQuery := query == "" || strings.Contains(strings.ToLower(app.Name), query) || strings.Contains(strings.ToLower(app.Description), query)
				matchesCategory := category == "" || strings.EqualFold(app.Category, category)
				if matchesQuery && matchesCategory {
					apps = append(apps, app)
				}
			}
			value = apps
		}
	}
	respond(w, value, err, http.StatusOK)
}

func (s *server) catalogSources(w http.ResponseWriter, r *http.Request) {
	value, err := s.provider.CatalogSources(r.Context())
	respond(w, value, err, http.StatusOK)
}

func (s *server) appAction(w http.ResponseWriter, r *http.Request) {
	action := r.PathValue("action")
	if action != "install" && action != "start" && action != "stop" && action != "restart" && action != "update" && action != "downgrade" {
		writeError(w, http.StatusBadRequest, "INVALID_APP_ACTION", "unsupported app action")
		return
	}
	id := r.PathValue("id")
	job := s.jobs.create("app."+action, id)
	s.auditStore.add("app."+action, id, "accepted", "job="+job.ID)
	s.eventHub.publish("job", job.ID+" accepted")
	go func() {
		ctx := context.WithoutCancel(r.Context())
		s.jobs.run(ctx, job, func(operationContext context.Context) (platform.Job, error) {
			if !safety.RealOperationsEnabled() {
				message := "app " + action + " planned; real operations are disabled"
				s.auditStore.add("app."+action, id, "planned", message)
				s.notificationStore.add("Ação planejada", id+": "+message, "apps", "info")
				s.eventHub.publish("job", message)
				return platform.Job{ID: job.ID, Status: "success", Message: message}, nil
			}
			result, err := s.appActionForManifest(operationContext, id, action)
			if err != nil {
				s.notificationStore.add("Falha no App", id+": "+err.Error(), "apps", "error")
				return platform.Job{}, err
			}
			s.notificationStore.add("App atualizado", result.Message, "apps", "success")
			s.eventHub.publish("job", result.Message)
			return result, nil
		})
	}()
	if snapshot, ok := s.jobs.get(job.ID); ok {
		respond(w, snapshot, nil, http.StatusAccepted)
		return
	}
	respond(w, map[string]any{"id": job.ID, "status": "queued", "message": "Aguardando execução"}, nil, http.StatusAccepted)
}

type declarativeAppInstaller interface {
	AppActionWithCompose(context.Context, string, string, string) (platform.Job, error)
}

func (s *server) appActionForManifest(ctx context.Context, id, action string) (platform.Job, error) {
	if action == "install" {
		for _, item := range s.declarative.list() {
			if item.ID == id {
				installer, ok := s.provider.(declarativeAppInstaller)
				if !ok {
					return platform.Job{}, platform.ErrUnsupported
				}
				return installer.AppActionWithCompose(ctx, id, action, item.Compose)
			}
		}
	}
	return s.provider.AppAction(ctx, id, action)
}

func (s *server) removeApp(w http.ResponseWriter, r *http.Request) {
	job := s.jobs.create("app.remove", r.PathValue("id"))
	s.auditStore.add("app.remove", r.PathValue("id"), "accepted", "job="+job.ID)
	go func() {
		ctx := context.WithoutCancel(r.Context())
		s.jobs.update(job.ID, func(item *managedJob) { item.Status = "running"; item.Progress = 10; item.Message = "Removendo app" })
		if !safety.RealOperationsEnabled() {
			message := "app removal planned; real operations are disabled"
			s.auditStore.add("app.remove", r.PathValue("id"), "planned", message)
			s.notificationStore.add("Desinstalação planejada", r.PathValue("id")+": "+message, "apps", "info")
			s.jobs.update(job.ID, func(item *managedJob) { item.Status = "success"; item.Progress = 100; item.Message = message })
			return
		}
		_, err := s.provider.RemoveApp(ctx, r.PathValue("id"))
		if err != nil {
			s.auditStore.add("app.remove", r.PathValue("id"), "error", err.Error())
			s.notificationStore.add("Falha ao remover app", r.PathValue("id")+": "+err.Error(), "apps", "error")
			s.jobs.update(job.ID, func(item *managedJob) {
				item.Status = "error"
				item.Progress = 100
				item.Message = "Remoção falhou"
				item.Error = err.Error()
			})
			return
		}
		s.notificationStore.add("App removido", r.PathValue("id"), "apps", "success")
		s.eventHub.publish("job", r.PathValue("id")+" removed")
		s.auditStore.add("app.remove", r.PathValue("id"), "success", "")
		s.jobs.update(job.ID, func(item *managedJob) { item.Status = "success"; item.Progress = 100; item.Message = "App removido" })
	}()
	if snapshot, ok := s.jobs.get(job.ID); ok {
		respond(w, snapshot, nil, http.StatusAccepted)
		return
	}
	respond(w, map[string]any{"id": job.ID, "status": "queued", "message": "Aguardando execução"}, nil, http.StatusAccepted)
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

func (s *server) appHealth(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	apps, err := s.provider.Apps(r.Context(), false)
	if err != nil {
		respond(w, map[string]any{"id": id, "status": "unknown", "healthy": false, "message": "app status unavailable"}, err, http.StatusOK)
		return
	}
	for _, app := range apps {
		if app.ID == id {
			health := app.Health
			if health == "" {
				health = app.Status
			}
			respond(w, map[string]any{"id": id, "status": health, "healthy": health == "running" || health == "healthy", "dependencies": app.Dependencies, "message": healthMessage(health)}, nil, http.StatusOK)
			return
		}
	}
	writeError(w, http.StatusNotFound, "APP_NOT_FOUND", "app not found")
}

func healthMessage(status string) string {
	switch status {
	case "running", "healthy":
		return "app is healthy"
	case "stopped":
		return "app is stopped"
	default:
		return "app health is unknown"
	}
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
	if value, ok := s.jobs.get(r.PathValue("id")); ok {
		respond(w, value, nil, http.StatusOK)
		return
	}
	value, err := s.provider.Job(r.Context(), r.PathValue("id"))
	respond(w, value, err, http.StatusOK)
}

func (s *server) jobsList(w http.ResponseWriter, _ *http.Request) {
	respond(w, s.jobs.list(), nil, http.StatusOK)
}
func (s *server) cancelJob(w http.ResponseWriter, r *http.Request) {
	if !s.jobs.cancel(r.PathValue("id")) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "job not found or already started")
		return
	}
	respond(w, map[string]any{"cancelled": true, "id": r.PathValue("id")}, nil, http.StatusAccepted)
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
