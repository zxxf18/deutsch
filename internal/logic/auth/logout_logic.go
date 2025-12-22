package auth

import (
	"context"

	"deutsch/internal/code"
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

func (l *LogoutLogic) Logout() (resp *types.LogoutResponse, err error) {
	// todo 需要把 jwt 给作废掉
	resp = &types.LogoutResponse{}
	resp.Base = *code.BaseSuccessResp()
	return
}
