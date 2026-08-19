package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"labsos/backend/providers/mock"
)

func TestPasswordChangeRequiresSessionAfterSetup(t *testing.T) {
	t.Setenv("LABSOS_ADMIN_PASSWORD", "correct-password")
	t.Setenv("LABSOS_ADMIN_PASSWORD_FILE", t.TempDir()+"/password")
	handler := New(mock.New())
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/password", strings.NewReader(`{"current":"correct-password","new":"new-password"}`))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

func TestPasswordHashUsesSaltedPBKDF2(t *testing.T) {
	hashed, err := hashPassword("correct-password")
	if err != nil || !strings.HasPrefix(hashed, "pbkdf2-sha256$") {
		t.Fatalf("hash = %q, err = %v", hashed, err)
	}
	if !passwordMatches("correct-password", hashed) || passwordMatches("wrong-password", hashed) {
		t.Fatal("PBKDF2 password verification failed")
	}
}

func TestSessionStoreRestoresUnexpiredSession(t *testing.T) {
	t.Setenv("LABSOS_SESSIONS_FILE", t.TempDir()+"/sessions.json")
	first := newSessionStore()
	first.mu.Lock()
	first.sessions["token"] = session{Token: "token", CreatedAt: time.Now().UTC(), LastSeen: time.Now().UTC(), UserAgent: "test"}
	first.persistLocked()
	first.mu.Unlock()
	second := newSessionStore()
	if !second.valid("token") {
		t.Fatal("session was not restored")
	}
}
