package invitecode

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/common"
	"deutsch/internal/svc"
	"deutsch/internal/types"
	"deutsch/model/gormdb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateLogic {
	return &GenerateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateLogic) Generate(req *types.GenerateInviteCodeRequest) (resp *types.GenerateInviteCodeResponse, err error) {
	resp = &types.GenerateInviteCodeResponse{}

	if req.Count < 1 || req.Count > 100 {
		return nil, code.NewCodeError(code.CodeInviteCountExceeded)
	}

	userID := common.GetUserID(l.ctx)
	if userID == "" {
		return nil, code.NewCodeError(code.CodeUnauthorized)
	}

	list, err := l.svcCtx.InviteCodeRepo.CreateBatch(l.ctx, req.Count, userID)
	if err != nil {
		l.Errorf("failed to create invite codes: %+v", err)
		return nil, code.NewCodeError(code.CodeInviteGenerationFailed)
	}

	resp.Base = *code.BaseSuccessResp()
	resp.Data.Items = make([]types.InviteCode, 0, len(list))
	for _, ic := range list {
		resp.Data.Items = append(resp.Data.Items, toTypesInviteCode(ic))
	}
	return resp, nil
}

func toTypesInviteCode(ic *gormdb.InviteCode) types.InviteCode {
	return types.InviteCode{
		Id:        ic.ID,
		Code:      ic.Code,
		UsedBy:    ic.UsedBy,
		ExpiresAt: ic.ExpiresAt.UnixMilli(),
		CreatedBy: ic.CreatedBy,
		IsEnabled: ic.IsEnabled,
		CreatedAt: ic.CreatedAt.UnixMilli(),
	}
}
