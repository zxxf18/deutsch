package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4/request"
	"github.com/zeromicro/go-zero/rest/httpx"

	"deutsch/internal/logic/auth"
	"deutsch/internal/svc"
)

func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, _ := request.AuthorizationHeaderExtractor.ExtractToken(r)
		l := auth.NewLogoutLogic(r.Context(), svcCtx)
		resp, err := l.Logout(token)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
