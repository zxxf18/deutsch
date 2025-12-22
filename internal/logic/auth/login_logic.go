package auth

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/pkg/jwt"
	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	jwtImpl jwt.JWT
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	l := &LoginLogic{
		Logger: logx.WithContext(logx.ContextWithFields(ctx, logx.Field("logic", "login"))),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
	jwtImpl, err := jwt.New(ctx, svcCtx.Config.JWTAuth.AccessSecret, svcCtx.Config.JWTAuth.AccessExpire)
	if err != nil {
		l.Errorf("failed to get jwt object: %+v", err)
	}
	l.jwtImpl = jwtImpl
	return l
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	resp = &types.LoginResponse{}
	jwtData := map[string]string{
		"userID": "zx",
	}

	jwtResp, err := l.jwtImpl.GenerateJWT(jwtData)
	if err != nil {
		l.Errorf("failed to generate jwt: %+v", err)
		return nil, err
	}
	resp.Base = *code.BaseSuccessResp()
	resp.Data.Expires = jwtResp.Expire.UnixMilli()
	resp.Data.MaxRefresh = jwtResp.MaxRefresh.UnixMilli()
	resp.Data.JwtToken = jwtResp.Token
	resp.Data.ID = "1001"
	resp.Data.Role = "admin"
	return
}
