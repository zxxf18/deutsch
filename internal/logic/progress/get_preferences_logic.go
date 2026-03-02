// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package progress

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/common"
	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPreferencesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPreferencesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPreferencesLogic {
	return &GetPreferencesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPreferencesLogic) GetPreferences() (resp *types.GetPreferencesResponse, err error) {
	userID := common.GetUserID(l.ctx)
	if userID == "" {
		return nil, code.NewCodeError(code.CodeUnauthorized)
	}

	pref, err := l.svcCtx.ProgressRepo.GetPreference(l.ctx, userID)
	if err != nil {
		l.Errorf("failed to get preference: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	resp = &types.GetPreferencesResponse{}
	resp.Base = *code.BaseSuccessResp()
	resp.Data.PreferredExamStateId = ""
	if pref != nil && pref.PreferredExamStateID != nil {
		resp.Data.PreferredExamStateId = *pref.PreferredExamStateID
	}
	return resp, nil
}
