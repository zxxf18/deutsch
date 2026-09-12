package middleware

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"deutsch/internal/sso"
	"deutsch/model/gormdb"
	"deutsch/model/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ensureUserFake embeds the production repository interface so this test stays
// compatible if unrelated repository methods are added.  Only the methods
// used by EnsureSSOUser are implemented in memory.
type ensureUserFake struct {
	repository.UserRepository

	mu              sync.Mutex
	users           map[string]*gormdb.User
	getErr          error
	createErr       error
	createCount     int
	failAfterCreate bool
}

func newEnsureUserFake() *ensureUserFake {
	return &ensureUserFake{users: make(map[string]*gormdb.User)}
}

func (f *ensureUserFake) GetByUserID(_ context.Context, id string) (*gormdb.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.getErr != nil {
		return nil, f.getErr
	}
	u, ok := f.users[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return u, nil
}

func (f *ensureUserFake) Create(_ context.Context, user *gormdb.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.createCount++
	if f.createErr != nil {
		if f.failAfterCreate {
			f.users[user.ID] = user
		}
		return f.createErr
	}
	if _, exists := f.users[user.ID]; exists {
		return errors.New("duplicate user id")
	}
	f.users[user.ID] = user
	return nil
}

func TestEnsureSSOUserCreatesByImmutableUUIDSubjectWithoutPassword(t *testing.T) {
	repo := newEnsureUserFake()
	sub := uuid.NewString()
	identity := sso.Identity{
		Subject:     sub,
		Email:       "new@example.com",
		Username:    "new-user",
		DisplayName: "New User",
		Role:        "user",
	}

	got, err := EnsureSSOUser(context.Background(), repo, identity)
	if err != nil {
		t.Fatalf("EnsureSSOUser() error = %v", err)
	}
	if got == nil {
		t.Fatal("EnsureSSOUser() returned nil user")
	}
	if got.ID != sub || got.Username != identity.Username || got.Email != identity.Email {
		t.Fatalf("created user identity mismatch: %+v", got)
	}
	if got.PasswordEncrypted != "" {
		t.Fatalf("SSO-provisioned user must not have a password, got %q", got.PasswordEncrypted)
	}
	if !got.IsEnabled || got.Role != "user" {
		t.Fatalf("unexpected default access state: enabled=%v role=%q", got.IsEnabled, got.Role)
	}
	if got.Nickname != identity.DisplayName {
		t.Fatalf("nickname = %q, want initial display name %q", got.Nickname, identity.DisplayName)
	}
}

func TestEnsureSSOUserRejectsNonUUIDSubject(t *testing.T) {
	repo := newEnsureUserFake()
	if got, err := EnsureSSOUser(context.Background(), repo, sso.Identity{
		Subject:  "casdoor-user-name",
		Email:    "user@example.com",
		Username: "user",
		Role:     "user",
	}); err == nil || got != nil {
		t.Fatalf("non-UUID subject should be rejected, got user=%+v err=%v", got, err)
	}
	if repo.createCount != 0 {
		t.Fatalf("non-UUID subject attempted creation %d times", repo.createCount)
	}
}

func TestEnsureSSOUserDoesNotOverwriteProfileOnRepeat(t *testing.T) {
	repo := newEnsureUserFake()
	sub := uuid.NewString()
	repo.users[sub] = &gormdb.User{
		ID:          sub,
		Username:    "stable-user",
		Email:       "stable@example.com",
		Nickname:    "手工昵称",
		Description: "手工简介",
		Role:        "user",
		IsEnabled:   true,
	}

	got, err := EnsureSSOUser(context.Background(), repo, sso.Identity{
		Subject:     sub,
		Email:       "changed@example.com",
		Username:    "changed-name",
		DisplayName: "Changed Display Name",
		Role:        "user",
	})
	if err != nil {
		t.Fatalf("EnsureSSOUser() error = %v", err)
	}
	if got.Nickname != "手工昵称" || got.Description != "手工简介" {
		t.Fatalf("existing profile was overwritten: nickname=%q description=%q", got.Nickname, got.Description)
	}
}

func TestEnsureSSOUserRejectsDisabledUser(t *testing.T) {
	repo := newEnsureUserFake()
	sub := uuid.NewString()
	repo.users[sub] = &gormdb.User{ID: sub, Username: "disabled", Email: "disabled@example.com", Role: "user", IsEnabled: false}

	if got, err := EnsureSSOUser(context.Background(), repo, sso.Identity{Subject: sub, Email: "disabled@example.com", Username: "disabled", Role: "user"}); err == nil || got != nil {
		t.Fatalf("disabled user should be rejected, got user=%+v err=%v", got, err)
	}
	if repo.createCount != 0 {
		t.Fatalf("disabled user attempted recreation %d times", repo.createCount)
	}
}

func TestEnsureSSOUserDoesNotRestoreSoftDeletedUser(t *testing.T) {
	repo := newEnsureUserFake()
	sub := uuid.NewString()
	repo.users[sub] = &gormdb.User{
		ID:        sub,
		Username:  "deleted",
		Email:     "deleted@example.com",
		Role:      "user",
		IsEnabled: true,
		DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
	}

	if got, err := EnsureSSOUser(context.Background(), repo, sso.Identity{Subject: sub, Email: "deleted@example.com", Username: "deleted", Role: "user"}); err == nil || got != nil {
		t.Fatalf("soft-deleted user should not be restored, got user=%+v err=%v", got, err)
	}
	if repo.createCount != 0 {
		t.Fatalf("soft-deleted user attempted recreation %d times", repo.createCount)
	}
}

func TestEnsureSSOUserRereadsAfterConcurrentCreateRace(t *testing.T) {
	repo := newEnsureUserFake()
	repo.createErr = errors.New("duplicate key")
	repo.failAfterCreate = true
	sub := uuid.NewString()

	got, err := EnsureSSOUser(context.Background(), repo, sso.Identity{
		Subject:     sub,
		Email:       "race@example.com",
		Username:    "race-user",
		DisplayName: "Race User",
		Role:        "user",
	})
	if err != nil {
		t.Fatalf("concurrent create race should reread winner: %v", err)
	}
	if got == nil || got.ID != sub {
		t.Fatalf("reread returned unexpected user: %+v", got)
	}
}

func TestEnsureSSOUserDoesNotBypassRepositoryFailure(t *testing.T) {
	repo := newEnsureUserFake()
	repo.getErr = errors.New("database unavailable")
	sub := uuid.NewString()

	if got, err := EnsureSSOUser(context.Background(), repo, sso.Identity{Subject: sub, Email: "db@example.com", Username: "db-user", Role: "user"}); err == nil || got != nil {
		t.Fatalf("database failure must not be bypassed, got user=%+v err=%v", got, err)
	}
	if repo.createCount != 0 {
		t.Fatalf("database failure attempted creation %d times", repo.createCount)
	}
}
