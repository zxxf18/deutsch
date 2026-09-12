package user

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"deutsch/internal/middleware"
	"deutsch/internal/sso"
	"deutsch/internal/svc"
	"deutsch/model/gormdb"
	"deutsch/model/repository"
	"gorm.io/gorm"
)

type profileMemoryStore struct {
	repository.UserRepository
	value *gormdb.User
}

func (s *profileMemoryStore) GetByUserID(_ context.Context, id string) (*gormdb.User, error) {
	if s.value == nil || s.value.ID != id {
		return nil, gorm.ErrRecordNotFound
	}
	copy := *s.value
	return &copy, nil
}
func (s *profileMemoryStore) Create(_ context.Context, u *gormdb.User) error { s.value = u; return nil }
func (s *profileMemoryStore) Update(_ context.Context, u *gormdb.User) error { s.value = u; return nil }

// Exercises real cookie verification, middleware provisioning, HTTP parsing,
// profile update and /me refresh; only storage is replaced with an isolated fake.
func TestSSOProfileHTTPProvisionSaveReload(t *testing.T) {
	const secret = "isolated-test-session-secret-not-for-production"
	const subject = "00000000-0000-4000-8000-000000000001"
	payload, _ := json.Marshal(map[string]any{"sub": subject, "email": "test@example.invalid", "username": "test-user", "display_name": "Initial", "role": "user", "email_verified": true, "exp": time.Now().Add(time.Hour).Unix()})
	encoded := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(encoded))
	cookie := &http.Cookie{Name: "test_session", Value: encoded + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))}
	repo := &profileMemoryStore{}
	auth := sso.New(sso.Config{CookieName: "test_session", SessionSecret: secret})
	sc := &svc.ServiceContext{SSO: auth, UserRepo: repo, SSOMiddleware: middleware.NewSSOMiddleware(auth, repo).Handle}
	request := func(method, path, body string, handler http.HandlerFunc) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, path, strings.NewReader(body))
		r.AddCookie(cookie)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, r)
		return w
	}
	first := request("GET", "/api/v1/auth/me", "", MeHandler(sc))
	if first.Code != 200 || repo.value == nil || repo.value.ID != subject {
		t.Fatalf("provision failed: %d %s", first.Code, first.Body.String())
	}
	if repo.value.PasswordEncrypted != "" {
		t.Fatal("SSO profile must not store a password")
	}
	saved := request("PATCH", "/api/v1/user/profile", `{"nickname":"修改后的昵称","description":"保存的简介"}`, sc.SSOMiddleware(UpdateUserHandler(sc)))
	if saved.Code != 200 || repo.value.Nickname != "修改后的昵称" || repo.value.Description != "保存的简介" {
		t.Fatalf("save failed: %d %s", saved.Code, saved.Body.String())
	}
	reloaded := request("GET", "/api/v1/auth/me", "", MeHandler(sc))
	var result map[string]any
	if err := json.Unmarshal(reloaded.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result["display_name"] != "修改后的昵称" || result["description"] != "保存的简介" {
		t.Fatalf("stale profile after reload: %s", reloaded.Body.String())
	}
	cleared := request("PATCH", "/api/v1/user/profile", `{"nickname":""}`, sc.SSOMiddleware(UpdateUserHandler(sc)))
	if cleared.Code != 200 || repo.value.Nickname != "" || repo.value.Description != "保存的简介" {
		t.Fatal("partial update or clear failed")
	}
	anonymous := httptest.NewRecorder()
	MeHandler(sc)(anonymous, httptest.NewRequest("GET", "/api/v1/auth/me", nil))
	if anonymous.Code != 401 {
		t.Fatalf("anonymous status: %d", anonymous.Code)
	}
}
