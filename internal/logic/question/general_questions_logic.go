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

type GeneralQuestionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGeneralQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GeneralQuestionsLogic {
	return &GeneralQuestionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GeneralQuestionsLogic) GeneralQuestions() (resp *types.ListQuestionsResponse, err error) {
	resp = &types.ListQuestionsResponse{}

	questions, err := l.svcCtx.QuestionRepo.GetGeneral(l.ctx)
	if err != nil {
		l.Errorf("failed to get general questions: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	states, _ := l.svcCtx.ConfigRepo.ListStates(l.ctx)
	idToSlug := buildStateIDToSlug(states)
	items, err := buildQuestionItems(l.ctx, l.svcCtx.QuestionRepo, questions, idToSlug)
	if err != nil {
		return nil, err
	}

	resp.Base = *code.BaseSuccessResp()
	resp.Data.Total = int64(len(items))
	resp.Data.Items = items
	return resp, nil
}
