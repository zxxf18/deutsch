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

type TrialCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTrialCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrialCheckLogic {
	return &TrialCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TrialCheckLogic) TrialCheck(req *types.TrialCheckRequest) (resp *types.TrialCheckResponse, err error) {
	resp = &types.TrialCheckResponse{}
	resp.Data.Results = make([]types.TrialCheckResultItem, 0)

	if len(req.Answers) == 0 {
		resp.Base = *code.BaseSuccessResp()
		return resp, nil
	}

	questionIDs := make([]string, 0, len(req.Answers))
	for id := range req.Answers {
		questionIDs = append(questionIDs, id)
	}

	optsMap, err := l.svcCtx.QuestionRepo.GetOptionsByQuestionIDs(l.ctx, questionIDs)
	if err != nil {
		l.Errorf("failed to get question options: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	for questionID, chosenAnswer := range req.Answers {
		optList := optsMap[questionID]
		correctOptionIndex := -1
		for _, o := range optList {
			if o.IsCorrect {
				correctOptionIndex = o.OptionIndex
				break
			}
		}
		correct := correctOptionIndex >= 0 && chosenAnswer == correctOptionIndex
		resp.Data.Results = append(resp.Data.Results, types.TrialCheckResultItem{
			QuestionId:         questionID,
			Correct:            correct,
			CorrectOptionIndex: correctOptionIndex,
		})
	}

	resp.Base = *code.BaseSuccessResp()
	return resp, nil
}
