package invitecode

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetInviteCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInviteCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInviteCodeLogic {
	return &GetInviteCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInviteCodeLogic) GetInviteCode(req *types.GetInviteCodeRequest) (resp *types.GetInviteCodeResponse, err error) {
	resp = &types.GetInviteCodeResponse{}

	ic, err := getInviteByIDOrCode(l.ctx, l.svcCtx.InviteCodeRepo, req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.NewCodeError(code.CodeInviteNotFound)
		}
		l.Errorf("failed to get invite code: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	resp.Base = *code.BaseSuccessResp()
	resp.Data = struct{ types.InviteCode }{InviteCode: toTypesInviteCode(ic)}
	return resp, nil
}
