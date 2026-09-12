package middleware

import (
	"errors"
	"net/http"

	"deutsch/internal/common"
	"deutsch/internal/sso"
	"deutsch/model/repository"
)

type SSOMiddleware struct {
	service *sso.Service
	users   repository.UserRepository
}

func NewSSOMiddleware(service *sso.Service, users repository.UserRepository) *SSOMiddleware {
	return &SSOMiddleware{service: service, users: users}
}

func (m *SSOMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return m.service.Require(func(w http.ResponseWriter, r *http.Request) {
		identity, ok := sso.IdentityFromContext(r.Context())
		if !ok || identity.Subject == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := EnsureSSOUser(r.Context(), m.users, identity)
		if err != nil {
			status := http.StatusServiceUnavailable
			if errors.Is(err, ErrUserDisabled) {
				status = http.StatusForbidden
			}
			http.Error(w, "business profile unavailable", status)
			return
		}
		r = r.WithContext(common.WithUserID(r.Context(), user.ID))
		r = r.WithContext(common.WithRole(r.Context(), identity.Role))
		next(w, r)
	})
}
