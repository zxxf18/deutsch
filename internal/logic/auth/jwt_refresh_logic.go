package auth

import (
	"context"
	"net/http"

	"deutsch/internal/code"
	"deutsch/internal/pkg/jwt"
	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type JwtRefreshLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJwtRefreshLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JwtRefreshLogic {
	return &JwtRefreshLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JwtRefreshLogic) JwtRefresh(r *http.Request) (resp *types.JwtRefreshResponse, err error) {
	resp = &types.JwtRefreshResponse{}

	jwtImpl, err := jwt.New(l.ctx, l.svcCtx.Config.JWTAuth.AccessSecret, l.svcCtx.Config.JWTAuth.AccessExpire)
	if err != nil {
		l.Errorf("failed to get jwt: %+v", err)
		return nil, code.NewCodeError(code.CodeInternalServerError)
	}

	jwtResp, err := jwtImpl.RefreshJWT(r)
	if err != nil {
		l.Errorf("failed to refresh jwt: %+v", err)
		return nil, code.NewCodeError(code.CodeInvalidToken)
	}

	resp.Base = *code.BaseSuccessResp()
	resp.Data.JwtToken = jwtResp.Token
	resp.Data.Expires = jwtResp.Expire.UnixMilli()
	resp.Data.MaxRefresh = jwtResp.MaxRefresh.UnixMilli()
	return
}
