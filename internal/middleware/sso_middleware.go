package middleware

import (
	"net/http"

	"deutsch/internal/common"
	"deutsch/internal/sso"
)

type SSOMiddleware struct{ service *sso.Service }

func NewSSOMiddleware(service *sso.Service) *SSOMiddleware { return &SSOMiddleware{service: service} }

func (m *SSOMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return m.service.Require(func(w http.ResponseWriter, r *http.Request) {
		subject, role, ok := sso.SessionFromContext(r.Context())
		if !ok || subject == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		r = r.WithContext(common.WithUserID(r.Context(), subject))
		r = r.WithContext(common.WithRole(r.Context(), role))
		next(w, r)
	})
}
