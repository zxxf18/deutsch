package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4/request"
	"github.com/zeromicro/go-zero/rest/httpx"

	"deutsch/internal/code"
	"deutsch/internal/pkg/blacklist"
	"deutsch/internal/types"
)

type BlacklistMiddleware struct {
	blacklist *blacklist.TokenBlacklist
}

func NewBlacklistMiddleware(bl *blacklist.TokenBlacklist) *BlacklistMiddleware {
	return &BlacklistMiddleware{blacklist: bl}
}

func (m *BlacklistMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := request.AuthorizationHeaderExtractor.ExtractToken(r)
		if err != nil || token == "" {
			next(w, r)
			return
		}
		ok, err := m.blacklist.IsBlacklisted(token)
		if err != nil {
			httpx.WriteJson(w, http.StatusInternalServerError, &types.Base{
				Code: int(code.CodeInternalServerError),
				Msg:  code.CodeInternalServerError.Message(),
			})
			return
		}
		if ok {
			httpx.WriteJson(w, http.StatusUnauthorized, &types.Base{
				Code: int(code.CodeTokenBlacklisted),
				Msg:  code.CodeTokenBlacklisted.Message(),
			})
			return
		}
		next(w, r)
	}
}
