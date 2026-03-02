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

// TrialQuestionCount 游客体验题目数量（与 config 保持一致）
const TrialQuestionCount = 10

type TrialQuestionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTrialQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrialQuestionsLogic {
	return &TrialQuestionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TrialQuestionsLogic) TrialQuestions() (resp *types.TrialListQuestionsResponse, err error) {
	resp = &types.TrialListQuestionsResponse{}

	generalList, err := l.svcCtx.QuestionRepo.GetRandomGeneral(l.ctx, TrialQuestionCount)
	if err != nil {
		l.Errorf("failed to get trial questions: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	states, _ := l.svcCtx.ConfigRepo.ListStates(l.ctx)
	idToSlug := BuildStateIDToSlug(states)
	items, err := BuildTrialQuestionItems(l.ctx, l.svcCtx.QuestionRepo, generalList, idToSlug)
	if err != nil {
		return nil, err
	}

	resp.Base = *code.BaseSuccessResp()
	resp.Data.Total = int64(len(items))
	resp.Data.Items = items
	return resp, nil
}
