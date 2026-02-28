package invitecode

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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
	resp = &types.EnableInviteCodeResponse{}

	ic, err := getInviteByIDOrCode(l.ctx, l.svcCtx.InviteCodeRepo, req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.NewCodeError(code.CodeInviteNotFound)
		}
		l.Errorf("failed to get invite code: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	// 已使用的邀请码不允许修改启用状态
	if ic.UsedBy != "" {
		return nil, code.NewCodeError(code.CodeInviteCodeAlreadyUsed)
	}

	ic.IsEnabled = req.IsEnabled
	if err := l.svcCtx.InviteCodeRepo.Update(l.ctx, ic); err != nil {
		l.Errorf("failed to update invite code: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	resp.Base = *code.BaseSuccessResp()
	return resp, nil
}
