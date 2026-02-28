package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"deutsch/internal/code"
	"deutsch/internal/common"
	"deutsch/internal/config"
	"deutsch/internal/types"
)

type AdminMiddleware struct {
	cfg config.Config
}

func NewAdminMiddleware(c config.Config) *AdminMiddleware {
	return &AdminMiddleware{cfg: c}
}

func (m *AdminMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := common.GetRole(r.Context())
		if role != "admin" {
			httpx.WriteJson(w, http.StatusForbidden, &types.Base{
				Code: int(code.CodeAdminRequired),
				Msg:  code.CodeAdminRequired.Message(),
			})
			return
		}
		next(w, r)
	}
}
