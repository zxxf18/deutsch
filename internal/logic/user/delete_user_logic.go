package user

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/common"
	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserRequest) (resp *types.DeleteUserResponse, err error) {
	resp = &types.DeleteUserResponse{}

	currentUserID := common.GetUserID(l.ctx)
	if currentUserID == req.ID {
		return nil, code.NewCodeError(code.CodeValidationError)
	}

	user, err := l.svcCtx.UserRepo.GetByUserID(l.ctx, req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.NewCodeError(code.CodeUserNotFound)
		}
		l.Errorf("failed to get user: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	if err := l.svcCtx.UserRepo.Delete(l.ctx, user); err != nil {
		l.Errorf("failed to delete user: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	resp.Base = *code.BaseSuccessResp()
	return
}
