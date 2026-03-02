// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package question

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StateQuestionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStateQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StateQuestionsLogic {
	return &StateQuestionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StateQuestionsLogic) StateQuestions(req *types.StateQuestionsRequest) (resp *types.ListQuestionsResponse, err error) {
	resp = &types.ListQuestionsResponse{}

	if req.StateId == "" {
		return nil, code.NewCodeErrorWithMsg(code.CodeStateNotFound, "州ID不能为空，请检查")
	}

	states, err := l.svcCtx.ConfigRepo.ListStates(l.ctx)
	if err != nil {
		l.Errorf("failed to list states: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}
	stateExists := false
	for _, s := range states {
		if s.ID == req.StateId {
			stateExists = true
			break
		}
	}
	if !stateExists {
		return nil, code.NewCodeErrorWithMsg(code.CodeStateNotFound, "州ID不存在，请检查")
	}

	questions, err := l.svcCtx.QuestionRepo.GetByStateID(l.ctx, req.StateId)
	if err != nil {
		l.Errorf("failed to get state questions: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	idToSlug := BuildStateIDToSlug(states)
	items, err := BuildQuestionItems(l.ctx, l.svcCtx.QuestionRepo, questions, idToSlug)
	if err != nil {
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	resp.Base = *code.BaseSuccessResp()
	resp.Data.Total = int64(len(items))
	resp.Data.Items = items
	return resp, nil
}
