package invitecode

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateLogic {
	return &ValidateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidateLogic) Validate(inviteCode string) (resp *types.Base, err error) {
	resp = &types.Base{}
	if inviteCode == "" {
		return nil, code.NewCodeError(code.CodeInvalidInviteCode)
	}
	valid, _, err := l.svcCtx.InviteCodeRepo.Validate(l.ctx, inviteCode)
	if err != nil {
		l.Errorf("failed to validate invite code: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}
	if !valid {
		return nil, code.NewCodeError(code.CodeInvalidInviteCode)
	}
	*resp = *code.BaseSuccessResp()
	return resp, nil
}
