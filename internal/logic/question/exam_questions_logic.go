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

type ExamQuestionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExamQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExamQuestionsLogic {
	return &ExamQuestionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExamQuestionsLogic) ExamQuestions(req *types.GetExamQuestionsRequest) (resp *types.ListQuestionsResponse, err error) {
	resp = &types.ListQuestionsResponse{}

	// 30 通用 + 3 州，共 33 题
	const examGeneral = 30
	const examState = 3

	generalList, err := l.svcCtx.QuestionRepo.GetRandomGeneral(l.ctx, examGeneral)
	if err != nil {
		l.Errorf("failed to get exam general questions: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	var stateList []*gormdb.Question
	if req.StateId != "" {
		stateList, err = l.svcCtx.QuestionRepo.GetRandomByStateID(l.ctx, req.StateId, examState)
		if err != nil {
			l.Errorf("failed to get exam state questions: %+v", err)
			return nil, code.NewCodeError(code.CodeDatabaseError)
		}
	}

	all := append([]*gormdb.Question(nil), generalList...)
	if len(stateList) > 0 {
		all = append(all, stateList...)
	}
	// 打乱顺序
	randShuffle(all)

	states, _ := l.svcCtx.ConfigRepo.ListStates(l.ctx)
	idToSlug := buildStateIDToSlug(states)
	items, err := buildQuestionItems(l.ctx, l.svcCtx.QuestionRepo, all, idToSlug)
	if err != nil {
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	resp.Base = *code.BaseSuccessResp()
	resp.Data.Total = int64(len(items))
	resp.Data.Items = items
	return resp, nil
}
