package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"labsos/backend/internal/api"
	"labsos/backend/providers/mock"
)

func TestSecretWriteIsPlannedAndDoesNotReturnValue(t *testing.T) {
	t.Setenv("LABSOS_SECRETS_DIR", t.TempDir())
	t.Setenv("LABSOS_ENABLE_REAL_OPERATIONS", "")
	t.Setenv("LABSOS_CONFIRM_REAL_OPERATIONS", "")
	handler := api.New(mock.New())
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPut, "/api/v1/secrets/registry-token", strings.NewReader(`{"value":"do-not-return"}`)))
	if response.Code != http.StatusAccepted || strings.Contains(response.Body.String(), "do-not-return") {
		t.Fatalf("response: %d %s", response.Code, response.Body.String())
	}
}

func TestSecretNameRejectsTraversal(t *testing.T) {
	handler := api.New(mock.New())
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPut, "/api/v1/secrets/..%2Fescape", strings.NewReader(`{"value":"x"}`)))
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}
