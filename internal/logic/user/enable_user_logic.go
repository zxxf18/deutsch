package user

import (
	"context"

	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EnableUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnableUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnableUserLogic {
	return &EnableUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnableUserLogic) EnableUser(req *types.EnableUserRequest) (resp *types.EnableUserResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
