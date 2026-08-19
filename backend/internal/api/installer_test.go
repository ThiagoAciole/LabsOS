package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"labsos/backend/internal/api"
	"labsos/backend/providers/mock"
)

func TestInstallerMockFlow(t *testing.T) {
	handler := api.New(mock.New())
	disks := httptest.NewRecorder()
	handler.ServeHTTP(disks, httptest.NewRequest(http.MethodGet, "/api/v1/installer/disks", nil))
	if disks.Code != http.StatusOK || !strings.Contains(disks.Body.String(), "/dev/sda") {
		t.Fatalf("disks: %d %s", disks.Code, disks.Body.String())
	}
	start := httptest.NewRecorder()
	handler.ServeHTTP(start, httptest.NewRequest(http.MethodPost, "/api/v1/installer/start", strings.NewReader(`{"operation":"fresh","disk":"sda","serverName":"LabsOS","password":"test-password"}`)))
	if start.Code != http.StatusAccepted {
		t.Fatalf("start: %d %s", start.Code, start.Body.String())
	}
	var job struct {
		ID        string `json:"id"`
		Simulated bool   `json:"simulated"`
		Password  string `json:"password"`
	}
	if err := json.NewDecoder(start.Body).Decode(&job); err != nil || job.ID == "" {
		t.Fatalf("job: %+v %v", job, err)
	}
	if !job.Simulated || job.Password != "" {
		t.Fatalf("installer safety contract leaked or misreported: %+v", job)
	}
	time.Sleep(1500 * time.Millisecond)
	status := httptest.NewRecorder()
	handler.ServeHTTP(status, httptest.NewRequest(http.MethodGet, "/api/v1/installer/jobs/"+job.ID, nil))
	if status.Code != http.StatusOK || !strings.Contains(status.Body.String(), `"status":"complete"`) {
		t.Fatalf("status: %d %s", status.Code, status.Body.String())
	}
}

func TestInstallerRebootIsPlannedByDefault(t *testing.T) {
	t.Setenv("LABSOS_ENABLE_REAL_OPERATIONS", "")
	t.Setenv("LABSOS_CONFIRM_REAL_OPERATIONS", "")
	handler := api.New(mock.New())
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/v1/installer/reboot", nil))
	if response.Code != http.StatusAccepted || !strings.Contains(response.Body.String(), `"planned":true`) {
		t.Fatalf("response: %d %s", response.Code, response.Body.String())
	}
}
