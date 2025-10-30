// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package invitecode

import (
	"context"

	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EnableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnableLogic {
	return &EnableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnableLogic) Enable(req *types.EnableInviteCodeRequest) (resp *types.EnableInviteCodeResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
