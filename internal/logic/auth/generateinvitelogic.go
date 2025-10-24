// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateInviteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateInviteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateInviteLogic {
	return &GenerateInviteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateInviteLogic) GenerateInvite(req *types.GenerateInviteRequest) (resp *types.GenerateInviteResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
