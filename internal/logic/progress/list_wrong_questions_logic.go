// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package progress

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/common"
	"deutsch/internal/logic/question"
	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWrongQuestionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListWrongQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWrongQuestionsLogic {
	return &ListWrongQuestionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListWrongQuestionsLogic) ListWrongQuestions(req *types.Filter) (resp *types.ListWrongQuestionsResponse, err error) {
	userID := common.GetUserID(l.ctx)
	if userID == "" {
		return nil, code.NewCodeError(code.CodeUnauthorized)
	}

	pageNo, pageSize := req.PageNo, req.PageSize
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (pageNo - 1) * pageSize

	ids, total, err := l.svcCtx.ProgressRepo.GetWrongQuestionIDs(l.ctx, userID, offset, pageSize)
	if err != nil {
		l.Errorf("failed to list wrong questions: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	resp = &types.ListWrongQuestionsResponse{}
	resp.Base = *code.BaseSuccessResp()
	resp.Data.Total = total
	resp.Data.Items = []types.QuestionItem{}

	if len(ids) == 0 {
		return resp, nil
	}

	questions, err := l.svcCtx.QuestionRepo.GetByIDs(l.ctx, ids)
	if err != nil {
		l.Errorf("failed to get questions: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	states, _ := l.svcCtx.ConfigRepo.ListStates(l.ctx)
	idToSlug := question.BuildStateIDToSlug(states)
	items, err := question.BuildQuestionItems(l.ctx, l.svcCtx.QuestionRepo, questions, idToSlug)
	if err != nil {
		l.Errorf("failed to build question items: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	idToItem := make(map[string]types.QuestionItem)
	for _, it := range items {
		idToItem[it.Id] = it
	}
	ordered := make([]types.QuestionItem, 0, len(ids))
	for _, id := range ids {
		if it, ok := idToItem[id]; ok {
			ordered = append(ordered, it)
		}
	}

	resp.Data.Items = ordered
	return resp, nil
}
