package api

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"labsos/backend/internal/safety"
)

type session struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	LastSeen  time.Time `json:"lastSeen"`
	UserAgent string    `json:"userAgent"`
}
type sessionStore struct {
	mu       sync.Mutex
	sessions map[string]session
	attempts map[string][]time.Time
	path     string
}

func newSessionStore() *sessionStore {
	path := os.Getenv("LABSOS_SESSIONS_FILE")
	if path == "" {
		path = "/tmp/labsos-sessions.json"
	}
	store := &sessionStore{sessions: map[string]session{}, attempts: map[string][]time.Time{}, path: path}
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &store.sessions)
	}
	now := time.Now().UTC()
	for token, item := range store.sessions {
		if token == "" || now.Sub(item.LastSeen) > 24*time.Hour {
			delete(store.sessions, token)
		}
	}
	return store
}

func (s *sessionStore) persistLocked() {
	data, err := json.MarshalIndent(s.sessions, "", "  ")
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return
	}
	tmp, err := os.CreateTemp(filepath.Dir(s.path), ".sessions-*.tmp")
	if err != nil {
		return
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err = tmp.Write(data); err == nil {
		err = tmp.Chmod(0600)
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err == nil {
		_ = os.Rename(tmpPath, s.path)
	}
}

func (s *server) authStatus(w http.ResponseWriter, r *http.Request) {
	token, _ := r.Cookie("labsos_session")
	configured := configuredPassword() != ""
	authenticated := token != nil && s.sessions.valid(token.Value)
	if !configured {
		authenticated = true
	}
	respond(w, map[string]any{"authenticated": authenticated, "configured": configured, "localDevelopmentBypass": !configured}, nil, http.StatusOK)
}
func (s *server) login(w http.ResponseWriter, r *http.Request) {
	if !s.sessions.allowLogin(r.RemoteAddr) {
		w.Header().Set("Retry-After", "60")
		writeError(w, http.StatusTooManyRequests, "RATE_LIMITED", "too many login attempts")
		return
	}
	var input struct {
		Password string `json:"password"`
	}
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "password is required")
		return
	}
	expected := configuredPassword()
	if expected == "" || !passwordMatches(input.Password, expected) {
		writeError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid credentials")
		return
	}
	token := randomToken()
	now := time.Now().UTC()
	s.sessions.mu.Lock()
	s.sessions.sessions[token] = session{Token: token, CreatedAt: now, LastSeen: now, UserAgent: r.UserAgent()}
	s.sessions.persistLocked()
	s.sessions.mu.Unlock()
	http.SetCookie(w, &http.Cookie{Name: "labsos_session", Value: token, Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode, Secure: r.TLS != nil, MaxAge: 86400})
	respond(w, map[string]any{"authenticated": true}, nil, http.StatusOK)
}

func (s *sessionStore) allowLogin(origin string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-time.Minute)
	recent := s.attempts[origin][:0]
	for _, item := range s.attempts[origin] {
		if item.After(cutoff) {
			recent = append(recent, item)
		}
	}
	allowed := len(recent) < 5
	s.attempts[origin] = append(recent, now)
	return allowed
}

func (s *server) sessionsList(w http.ResponseWriter, r *http.Request) {
	current := ""
	if cookie, err := r.Cookie("labsos_session"); err == nil {
		current = cookie.Value
	}
	s.sessions.mu.Lock()
	defer s.sessions.mu.Unlock()
	result := make([]map[string]any, 0, len(s.sessions.sessions))
	for token, item := range s.sessions.sessions {
		if time.Since(item.LastSeen) > 24*time.Hour {
			continue
		}
		id := token
		if len(id) > 12 {
			id = id[:12]
		}
		result = append(result, map[string]any{"id": id, "createdAt": item.CreatedAt, "lastSeen": item.LastSeen, "userAgent": item.UserAgent, "current": token == current})
	}
	respond(w, result, nil, http.StatusOK)
}

func (s *server) revokeSession(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "INVALID_SESSION", "session id is required")
		return
	}
	s.sessions.mu.Lock()
	defer s.sessions.mu.Unlock()
	for token := range s.sessions.sessions {
		if strings.HasPrefix(token, id) {
			delete(s.sessions.sessions, token)
			s.sessions.persistLocked()
			s.auditStore.add("auth.session.revoke", id, "success", "")
			respond(w, map[string]any{"revoked": true}, nil, http.StatusOK)
			return
		}
	}
	writeError(w, http.StatusNotFound, "NOT_FOUND", "session not found")
}

func (s *server) changePassword(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Current string `json:"current"`
		Next    string `json:"new"`
	}
	if json.NewDecoder(r.Body).Decode(&input) != nil || len(input.Next) < 8 {
		writeError(w, http.StatusBadRequest, "INVALID_PASSWORD", "new password must contain at least 8 characters")
		return
	}
	if !passwordMatches(input.Current, configuredPassword()) {
		writeError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "current password is invalid")
		return
	}
	if !safety.RealOperationsEnabled() {
		s.auditStore.add("auth.password", "admin", "planned", "real operations disabled")
		respond(w, map[string]any{"changed": false, "simulated": true, "message": "password change planned; real operations are disabled"}, nil, http.StatusAccepted)
		return
	}
	path := passwordFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		writeError(w, http.StatusInternalServerError, "PASSWORD_DIRECTORY_FAILED", err.Error())
		return
	}
	hashed, err := hashPassword(input.Next)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "PASSWORD_HASH_FAILED", err.Error())
		return
	}
	if err := os.WriteFile(path, []byte(hashed+"\n"), 0600); err != nil {
		writeError(w, http.StatusInternalServerError, "PASSWORD_WRITE_FAILED", err.Error())
		return
	}
	s.auditStore.add("auth.password", "admin", "success", "hashed password updated")
	respond(w, map[string]any{"changed": true, "simulated": false}, nil, http.StatusOK)
}

func passwordFilePath() string {
	if value := os.Getenv("LABSOS_ADMIN_PASSWORD_FILE"); value != "" {
		return value
	}
	return "/var/lib/labsos/admin-password.hash"
}
func configuredPassword() string {
	if data, err := os.ReadFile(passwordFilePath()); err == nil && strings.TrimSpace(string(data)) != "" {
		return strings.TrimSpace(string(data))
	}
	return os.Getenv("LABSOS_ADMIN_PASSWORD")
}
func passwordMatches(input, expected string) bool {
	if strings.HasPrefix(expected, "pbkdf2-sha256$") {
		parts := strings.Split(expected, "$")
		if len(parts) != 4 {
			return false
		}
		iterations, err := strconv.Atoi(parts[1])
		if err != nil || iterations < 100000 || iterations > 1000000 {
			return false
		}
		salt, err := base64.RawStdEncoding.DecodeString(parts[2])
		if err != nil {
			return false
		}
		want, err := base64.RawStdEncoding.DecodeString(parts[3])
		if err != nil {
			return false
		}
		got := pbkdf2SHA256([]byte(input), salt, iterations, len(want))
		return subtle.ConstantTimeCompare(got, want) == 1
	}
	if strings.HasPrefix(expected, "sha256:") {
		digest := sha256.Sum256([]byte(input))
		encoded := hex.EncodeToString(digest[:])
		return subtle.ConstantTimeCompare([]byte(encoded), []byte(strings.TrimPrefix(expected, "sha256:"))) == 1
	}
	return subtle.ConstantTimeCompare([]byte(input), []byte(expected)) == 1
}

func hashPassword(input string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	const iterations = 120000
	derived := pbkdf2SHA256([]byte(input), salt, iterations, 32)
	return "pbkdf2-sha256$" + strconv.Itoa(iterations) + "$" + base64.RawStdEncoding.EncodeToString(salt) + "$" + base64.RawStdEncoding.EncodeToString(derived), nil
}

func pbkdf2SHA256(password, salt []byte, iterations, length int) []byte {
	result := make([]byte, 0, length)
	for block := 1; len(result) < length; block++ {
		mac := hmac.New(sha256.New, password)
		mac.Write(salt)
		mac.Write([]byte{byte(block >> 24), byte(block >> 16), byte(block >> 8), byte(block)})
		u := mac.Sum(nil)
		t := append([]byte(nil), u...)
		for i := 1; i < iterations; i++ {
			mac = hmac.New(sha256.New, password)
			mac.Write(u)
			u = mac.Sum(nil)
			for index := range t {
				t[index] ^= u[index]
			}
		}
		result = append(result, t...)
	}
	return result[:length]
}
func (s *server) logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("labsos_session"); err == nil {
		s.sessions.mu.Lock()
		delete(s.sessions.sessions, cookie.Value)
		s.sessions.persistLocked()
		s.sessions.mu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{Name: "labsos_session", Value: "", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteStrictMode})
	respond(w, map[string]any{"authenticated": false}, nil, http.StatusOK)
}
func (s *sessionStore) valid(token string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.sessions[token]
	if !ok || time.Since(item.LastSeen) > 24*time.Hour {
		return false
	}
	item.LastSeen = time.Now().UTC()
	s.sessions[token] = item
	s.persistLocked()
	return true
}
func randomToken() string {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		digest := sha256.Sum256([]byte(time.Now().String()))
		return hex.EncodeToString(digest[:])
	}
	return hex.EncodeToString(raw)
}
