package invitecode

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.DeleteInviteCodeRequest) (resp *types.DeleteInviteCodeResponse, err error) {
	resp = &types.DeleteInviteCodeResponse{}

	ic, err := getInviteByIDOrCode(l.ctx, l.svcCtx.InviteCodeRepo, req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.NewCodeError(code.CodeInviteNotFound)
		}
		l.Errorf("failed to get invite code: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	// 已使用的邀请码不允许删除
	if ic.UsedBy != "" {
		return nil, code.NewCodeError(code.CodeInviteCodeAlreadyUsed)
	}

	if err := l.svcCtx.InviteCodeRepo.Delete(l.ctx, ic); err != nil {
		l.Errorf("failed to delete invite code: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	resp.Base = *code.BaseSuccessResp()
	return resp, nil
}
