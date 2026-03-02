// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package question

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/svc"
	"deutsch/internal/types"
	"deutsch/model/gormdb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AllQuestionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAllQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllQuestionsLogic {
	return &AllQuestionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllQuestionsLogic) AllQuestions() (resp *types.AllQuestionsByStateResponse, err error) {
	resp = &types.AllQuestionsByStateResponse{}
	resp.Data = make(map[string][]types.QuestionItem)

	questions, err := l.svcCtx.QuestionRepo.GetAll(l.ctx)
	if err != nil {
		l.Errorf("failed to get all questions: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	states, _ := l.svcCtx.ConfigRepo.ListStates(l.ctx)
	idToSlug := BuildStateIDToSlug(states)

	// 按 state key 分组：general 或州 slug
	grouped := make(map[string][]*gormdb.Question)
	for _, q := range questions {
		key := resolveState(q, idToSlug)
		grouped[key] = append(grouped[key], q)
	}

	for key, qs := range grouped {
		items, err := BuildQuestionItems(l.ctx, l.svcCtx.QuestionRepo, qs, idToSlug)
		if err != nil {
			return nil, code.NewCodeError(code.CodeDatabaseError)
		}
		resp.Data[key] = items
	}

	resp.Base = *code.BaseSuccessResp()
	return resp, nil
}
