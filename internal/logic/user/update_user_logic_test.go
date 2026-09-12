package user

import (
	"context"
	"strings"
	"testing"

	"deutsch/internal/common"
	"deutsch/internal/svc"
	"deutsch/internal/types"
	"deutsch/model/gormdb"
	"deutsch/model/repository"
)

type profileRepo struct {
	repository.UserRepository
	user    gormdb.User
	updates int
}

func (r *profileRepo) GetByUserID(context.Context, string) (*gormdb.User, error) {
	u := r.user
	return &u, nil
}
func (r *profileRepo) Update(_ context.Context, u *gormdb.User) error {
	r.user = *u
	r.updates++
	return nil
}

func TestUpdateProfileBoundariesAndPartialUpdates(t *testing.T) {
	for _, tc := range []struct {
		name                      string
		nickname                  *string
		description               *string
		wantName, wantDescription string
		wantErr                   bool
	}{
		{"clear nickname", stringPtr(""), nil, "", "saved bio", false},
		{"fifty Chinese characters", stringPtr(strings.Repeat("夜", 50)), nil, strings.Repeat("夜", 50), "saved bio", false},
		{"fifty one rejected", stringPtr(strings.Repeat("夜", 51)), nil, "saved", "saved bio", true},
		{"description only preserves nickname", nil, stringPtr("new bio"), "saved", "new bio", false},
		{"clear description", nil, stringPtr(""), "saved", "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &profileRepo{user: gormdb.User{ID: "00000000-0000-4000-8000-000000000001", Nickname: "saved", Description: "saved bio", IsEnabled: true}}
			logic := NewUpdateUserLogic(common.WithUserID(context.Background(), repo.user.ID), &svc.ServiceContext{UserRepo: repo})
			_, err := logic.UpdateUser(&types.UpdateUserRequest{Nickname: tc.nickname, Description: tc.description})
			if (err != nil) != tc.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if repo.user.Nickname != tc.wantName || repo.user.Description != tc.wantDescription {
				t.Fatalf("unexpected saved profile: %+v", repo.user)
			}
			if tc.wantErr && repo.updates != 0 {
				t.Fatal("invalid input changed profile")
			}
		})
	}
}
func stringPtr(v string) *string { return &v }
