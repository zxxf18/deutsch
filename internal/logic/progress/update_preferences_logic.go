// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package progress

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/common"
	"deutsch/internal/svc"
	"deutsch/internal/types"
	"deutsch/model/gormdb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePreferencesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePreferencesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePreferencesLogic {
	return &UpdatePreferencesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePreferencesLogic) UpdatePreferences(req *types.UpdatePreferencesRequest) (resp *types.UpdatePreferencesResponse, err error) {
	userID := common.GetUserID(l.ctx)
	if userID == "" {
		return nil, code.NewCodeError(code.CodeUnauthorized)
	}

	if req.PreferredExamStateId != "" {
		states, err := l.svcCtx.ConfigRepo.ListStates(l.ctx)
		if err != nil {
			l.Errorf("failed to list states: %+v", err)
			return nil, code.NewCodeError(code.CodeDatabaseError)
		}
		found := false
		for _, s := range states {
			if s.ID == req.PreferredExamStateId {
				found = true
				break
			}
		}
		if !found {
			return nil, code.NewCodeErrorWithMsg(code.CodeStateNotFound, "州ID不存在，请检查")
		}
	}

	pref := &gormdb.UserPreference{UserID: userID}
	if req.PreferredExamStateId != "" {
		pref.PreferredExamStateID = &req.PreferredExamStateId
	} else {
		pref.PreferredExamStateID = nil
	}

	if err := l.svcCtx.ProgressRepo.UpsertPreference(l.ctx, pref); err != nil {
		l.Errorf("failed to upsert preference: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	resp = &types.UpdatePreferencesResponse{}
	resp.Base = *code.BaseSuccessResp()
	return resp, nil
}
