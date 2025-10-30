// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"deutsch/internal/code"
	"deutsch/internal/config"
	"deutsch/internal/pkg/jwt"
	"deutsch/internal/types"
)

type JWTMiddleware struct {
	cfg config.Config
	ctx context.Context
	jwt jwt.JWT
}

func NewJWTMiddleware(cfg config.Config) *JWTMiddleware {
	ctx := context.Background()
	j, err := jwt.New(ctx, cfg.JWTAuth.AccessSecret, cfg.JWTAuth.AccessExpire)
	if err != nil {
		logx.Errorf("failed to init admin middleware, %+v", err)
	}
	return &JWTMiddleware{
		cfg: cfg,
		ctx: ctx,
		jwt: j,
	}
}

func (m *JWTMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := m.jwt.CheckAndParseJWT(r)
		if err != nil {
			httpx.WriteJson(w, http.StatusUnauthorized, &types.Base{
				Code: int(code.CodeInvalidToken),
				Msg:  code.CodeInvalidToken.Message(),
			})
			return
		}
		next(w, r)
	}
}
