package invitecode

import (
	"context"

	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInviteCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInviteCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInviteCodeLogic {
	return &GetInviteCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInviteCodeLogic) GetInviteCode() (resp *types.GetInviteCodeResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
