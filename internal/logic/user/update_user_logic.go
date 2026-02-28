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

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserRequest) (resp *types.UpdateUserResponse, err error) {
	resp = &types.UpdateUserResponse{}

	currentUserID := common.GetUserID(l.ctx)
	if currentUserID == "" {
		return nil, code.NewCodeError(code.CodeUnauthorized)
	}

	user, err := l.svcCtx.UserRepo.GetByUserID(l.ctx, currentUserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.NewCodeError(code.CodeUserNotFound)
		}
		l.Errorf("failed to get user: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	if req.Nickname != "" {
		if len(req.Nickname) > 50 {
			return nil, code.NewCodeError(code.CodeNicknameTooLong)
		}
		user.Nickname = req.Nickname
	}
	if len(req.Description) > 500 {
		return nil, code.NewCodeError(code.CodeDescriptionTooLong)
	}
	user.Description = req.Description

	if err := l.svcCtx.UserRepo.Update(l.ctx, user); err != nil {
		l.Errorf("failed to update user: %+v", err)
		return nil, code.NewCodeError(code.CodeProfileUpdateFailed)
	}

	resp.Base = *code.BaseSuccessResp()
	resp.Data.ID = user.ID
	resp.Data.Username = user.Username
	resp.Data.Email = user.Email
	if user.Phone != nil {
		resp.Data.Phone = *user.Phone
	}
	resp.Data.Role = user.Role
	resp.Data.Nickname = user.Nickname
	resp.Data.IsEnabled = user.IsEnabled
	resp.Data.Description = user.Description
	resp.Data.CreatedAt = user.CreatedAt.UnixMilli()
	resp.Data.UpdatedAt = user.UpdatedAt.UnixMilli()
	return
}
