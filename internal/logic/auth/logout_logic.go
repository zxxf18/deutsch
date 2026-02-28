package auth

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/pkg/jwt"
	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(tokenString string) (resp *types.LogoutResponse, err error) {
	resp = &types.LogoutResponse{}
	resp.Base = *code.BaseSuccessResp()

	if tokenString == "" {
		return resp, nil
	}
	expireAt, err := jwt.GetExpireFromToken(tokenString)
	if err != nil {
		l.Infof("logout: failed to parse token expire: %v", err)
		return resp, nil
	}
	if err := l.svcCtx.TokenBlacklist.AddWithExpire(tokenString, expireAt); err != nil {
		l.Errorf("logout: failed to add token to blacklist: %v", err)
		// 不影响响应，客户端可丢弃 token
	}
	return resp, nil
}
