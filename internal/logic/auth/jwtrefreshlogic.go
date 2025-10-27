// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type JWTRefreshLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJWTRefreshLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JWTRefreshLogic {
	return &JWTRefreshLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JWTRefreshLogic) JWTRefresh() (resp *types.JWTRefreshResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
