package user

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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
	resp = &types.EnableUserResponse{}

	user, err := l.svcCtx.UserRepo.GetByUserID(l.ctx, req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.NewCodeError(code.CodeUserNotFound)
		}
		l.Errorf("failed to get user: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	user.IsEnabled = req.IsEnabled
	if err := l.svcCtx.UserRepo.Update(l.ctx, user); err != nil {
		l.Errorf("failed to update user: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	resp.Base = *code.BaseSuccessResp()
	return
}
