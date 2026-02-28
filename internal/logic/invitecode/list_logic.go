package invitecode

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/svc"
	"deutsch/internal/types"
	"deutsch/model/repository"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.ListInviteCodeRequest) (resp *types.ListInviteCodeResponse, err error) {
	resp = &types.ListInviteCodeResponse{}

	filter := &repository.InviteCodeListFilter{
		Filter: repository.Filter{
			PageNo:   req.PageNo,
			PageSize: req.PageSize,
		},
		AvailableOnly: req.AvailableOnly,
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

	invites, total, err := l.svcCtx.InviteCodeRepo.List(l.ctx, filter)
	if err != nil {
		l.Errorf("failed to list invite codes: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	resp.Base = *code.BaseSuccessResp()
	resp.Data.PageNo = filter.PageNo
	resp.Data.PageSize = filter.PageSize
	resp.Data.Total = total
	resp.Data.Items = make([]types.InviteCode, 0, len(invites))
	for _, ic := range invites {
		resp.Data.Items = append(resp.Data.Items, toTypesInviteCode(ic))
	}
	return resp, nil
}
