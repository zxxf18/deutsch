package user

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/svc"
	"deutsch/internal/types"
	"deutsch/model/repository"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserLogic {
	return &ListUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUserLogic) ListUser(req *types.Filter) (resp *types.ListUserResponse, err error) {
	resp = &types.ListUserResponse{}

	filter := &repository.Filter{
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
	}
	if filter.PageNo <= 0 {
		filter.PageNo = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	users, total, err := l.svcCtx.UserRepo.List(l.ctx, filter)
	if err != nil {
		l.Errorf("failed to list users: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	resp.Base = *code.BaseSuccessResp()
	resp.Data.PageNo = filter.PageNo
	resp.Data.PageSize = filter.PageSize
	resp.Data.Total = total
	resp.Data.Items = make([]types.User, 0, len(users))
	for _, u := range users {
		var phone string
		if u.Phone != nil {
			phone = *u.Phone
		}
		resp.Data.Items = append(resp.Data.Items, types.User{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			Phone:     phone,
			Role:      u.Role,
			Nickname:  u.Nickname,
			IsEnabled: u.IsEnabled,
			CreatedAt: u.CreatedAt.UnixMilli(),
			UpdatedAt: u.UpdatedAt.UnixMilli(),
		})
	}
	return
}
