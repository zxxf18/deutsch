// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package invitecode

import (
	"context"

	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete() (resp *types.DeleteInviteCodeResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
