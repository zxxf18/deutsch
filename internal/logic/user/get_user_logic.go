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

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserRequest) (resp *types.GetUserResponse, err error) {
	resp = &types.GetUserResponse{}

	currentUserID := common.GetUserID(l.ctx)
	role := common.GetRole(l.ctx)

	if role != "admin" && currentUserID != req.ID {
		return nil, code.NewCodeError(code.CodeUnauthorized)
	}

	user, err := l.svcCtx.UserRepo.GetByUserID(l.ctx, req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.NewCodeError(code.CodeUserNotFound)
		}
		l.Errorf("failed to get user: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
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
