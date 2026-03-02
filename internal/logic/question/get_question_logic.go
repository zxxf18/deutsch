// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package question

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionLogic {
	return &GetQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQuestionLogic) GetQuestion(req *types.GetQuestionRequest) (resp *types.GetQuestionResponse, err error) {
	resp = &types.GetQuestionResponse{}

	if req.QuestionId == "" {
		return nil, code.NewCodeError(code.CodeQuestionNotFound)
	}

	q, err := l.svcCtx.QuestionRepo.GetByID(l.ctx, req.QuestionId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.NewCodeError(code.CodeQuestionNotFound)
		}
		l.Errorf("failed to get question: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	optsMap, err := l.svcCtx.QuestionRepo.GetOptionsByQuestionIDs(l.ctx, []string{q.ID})
	if err != nil {
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}
	opts := optsMap[q.ID]

	state := "general"
	if q.StateID != nil && *q.StateID != "" {
		states, _ := l.svcCtx.ConfigRepo.ListStates(l.ctx)
		idToSlug := buildStateIDToSlug(states)
		state = resolveState(q, idToSlug)
	}

	resp.Base = *code.BaseSuccessResp()
	resp.Data = toQuestionItem(q, opts, state)
	return resp, nil
}
