package auth

import (
	"context"

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

func (l *JwtRefreshLogic) JwtRefresh() (resp *types.JwtRefreshResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
