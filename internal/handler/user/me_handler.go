package user

import (
	"net/http"

	"deutsch/internal/common"
	"deutsch/internal/sso"
	"deutsch/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Me reads the saved business profile, not the old display name in the SSO cookie.
func MeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return svcCtx.SSOMiddleware(func(w http.ResponseWriter, r *http.Request) {
		user, err := svcCtx.UserRepo.GetByUserID(r.Context(), common.GetUserID(r.Context()))
		if err != nil {
			http.Error(w, "business profile unavailable", http.StatusServiceUnavailable)
			return
		}
		identity, _ := sso.IdentityFromContext(r.Context())
		w.Header().Set("Cache-Control", "no-store")
		httpx.OkJsonCtx(r.Context(), w, map[string]any{
			"sub": user.ID, "username": identity.Username, "email": identity.Email,
			"display_name": user.Nickname, "description": user.Description,
			"role": user.Role, "email_verified": true,
		})
	})
}
