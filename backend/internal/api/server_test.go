package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"labsos/backend/internal/api"
	"labsos/backend/internal/platform"
	"labsos/backend/providers/mock"
)

func TestReadyRoutesReturnJSON(t *testing.T) {
	handler := api.New(mock.New())
	for _, path := range []string{
		"/api/v1/system/summary",
		"/api/v1/system/health",
		"/api/v1/apps",
		"/api/v1/catalog/apps",
		"/api/v1/settings",
		"/api/v1/events",
		"/api/v1/scheduler/tasks",
	} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		if response.Code != http.StatusOK {
			t.Fatalf("GET %s returned %d: %s", path, response.Code, response.Body.String())
		}
		if got := response.Header().Get("Content-Type"); got != "application/json" {
			t.Fatalf("GET %s content type = %q", path, got)
		}
	}
}

func TestAppActionCreatesRetrievableJob(t *testing.T) {
	t.Setenv("LABSOS_ENABLE_REAL_OPERATIONS", "false")
	t.Setenv("LABSOS_CONFIRM_REAL_OPERATIONS", "")
	handler := api.New(mock.New())
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/v1/apps/jellyfin/stop", nil))
	if response.Code != http.StatusAccepted {
		t.Fatalf("action returned %d: %s", response.Code, response.Body.String())
	}
	var job platform.Job
	if err := json.NewDecoder(response.Body).Decode(&job); err != nil {
		t.Fatal(err)
	}
	lookup := httptest.NewRecorder()
	handler.ServeHTTP(lookup, httptest.NewRequest(http.MethodGet, "/api/v1/jobs/"+job.ID, nil))
	if lookup.Code != http.StatusOK {
		t.Fatalf("job lookup returned %d: %s", lookup.Code, lookup.Body.String())
	}
}

func TestAppActionIsPlannedWithoutRealOperations(t *testing.T) {
	t.Setenv("LABSOS_ENABLE_REAL_OPERATIONS", "false")
	t.Setenv("LABSOS_CONFIRM_REAL_OPERATIONS", "")
	handler := api.New(mock.New())
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/v1/apps/jellyfin/stop", nil))
	if response.Code != http.StatusAccepted || !strings.Contains(response.Body.String(), "job-") {
		t.Fatalf("action response: %d %s", response.Code, response.Body.String())
	}
	var job platform.Job
	if err := json.NewDecoder(response.Body).Decode(&job); err != nil {
		t.Fatal(err)
	}
	for range 20 {
		lookup := httptest.NewRecorder()
		handler.ServeHTTP(lookup, httptest.NewRequest(http.MethodGet, "/api/v1/jobs/"+job.ID, nil))
		if strings.Contains(lookup.Body.String(), "planned") {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("planned job did not finish")
}

func TestSettingsUpdateAndPowerSimulation(t *testing.T) {
	handler := api.New(mock.New())
	settings := httptest.NewRecorder()
	handler.ServeHTTP(settings, httptest.NewRequest(http.MethodPut, "/api/v1/settings/system", strings.NewReader(`{"hostname":"lab"}`)))
	if settings.Code != http.StatusOK || !strings.Contains(settings.Body.String(), `"hostname":"lab"`) {
		t.Fatalf("settings response: %d %s", settings.Code, settings.Body.String())
	}
	power := httptest.NewRecorder()
	handler.ServeHTTP(power, httptest.NewRequest(http.MethodPost, "/api/v1/system/reboot", nil))
	if power.Code != http.StatusAccepted || !strings.Contains(power.Body.String(), "simulated") {
		t.Fatalf("power response: %d %s", power.Code, power.Body.String())
	}
}

func TestUpdateIsPlannedWithoutRealOperations(t *testing.T) {
	t.Setenv("LABSOS_ENABLE_REAL_OPERATIONS", "false")
	t.Setenv("LABSOS_CONFIRM_REAL_OPERATIONS", "")
	handler := api.New(mock.New())
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/v1/system/update", nil))
	if response.Code != http.StatusAccepted || !strings.Contains(response.Body.String(), `"simulated":true`) {
		t.Fatalf("update response: %d %s", response.Code, response.Body.String())
	}
}

func TestErrorsAreStructured(t *testing.T) {
	handler := api.New(mock.New())
	for _, request := range []*http.Request{
		httptest.NewRequest(http.MethodPost, "/api/v1/settings", nil),
		httptest.NewRequest(http.MethodGet, "/api/v1/jobs/missing", nil),
		httptest.NewRequest(http.MethodGet, "/api/v1/unknown", nil),
	} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code < 400 || !strings.Contains(response.Body.String(), `"error"`) {
			t.Fatalf("%s %s returned %d %s", request.Method, request.URL.Path, response.Code, response.Body.String())
		}
	}
}

func TestAppLogsEndpointReturnsRecentLogs(t *testing.T) {
	handler := api.New(mock.New())
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/apps/jellyfin/logs?lines=20", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
}
