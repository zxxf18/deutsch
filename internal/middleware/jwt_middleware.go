package middleware

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"deutsch/internal/code"
	"deutsch/internal/common"
	"deutsch/internal/config"
	"deutsch/internal/pkg/jwt"
	"deutsch/internal/types"
)

const jwtKey = "jwt"

type JWTMiddleware struct {
	cfg config.Config
	ctx context.Context
	jwt jwt.JWT
}

func NewJWTMiddleware(cfg config.Config) *JWTMiddleware {
	ctx := context.Background()
	j, err := jwt.New(ctx, cfg.JWTAuth.AccessSecret, cfg.JWTAuth.AccessExpire)
	if err != nil {
		logx.Errorf("failed to init jwt middleware, %+v", err)
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
		jwtData, _ := data[jwtKey].(map[string]any)
		if jwtData != nil {
			if userID, ok := jwtData["userID"].(string); ok {
				r = r.WithContext(common.WithUserID(r.Context(), userID))
			}
			if role, ok := jwtData["role"].(string); ok {
				r = r.WithContext(common.WithRole(r.Context(), role))
			}
		}
		next(w, r)
	}
}
