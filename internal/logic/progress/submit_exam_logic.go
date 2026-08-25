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

type SubmitExamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitExamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitExamLogic {
	return &SubmitExamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

const passScore = 17

func (l *SubmitExamLogic) SubmitExam(req *types.SubmitExamRequest) (resp *types.SubmitExamResponse, err error) {
	userID := common.GetUserID(l.ctx)
	if userID == "" {
		return nil, code.NewCodeError(code.CodeUnauthorized)
	}

	if len(req.Answers) == 0 {
		return nil, code.NewCodeErrorWithMsg(code.CodeAnswersInvalid, "答案不能为空")
	}

	questionIDs := make([]string, 0, len(req.Answers))
	for qid := range req.Answers {
		questionIDs = append(questionIDs, qid)
	}

	optsMap, err := l.svcCtx.QuestionRepo.GetOptionsByQuestionIDs(l.ctx, questionIDs)
	if err != nil {
		l.Errorf("failed to get options: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	questions, err := l.svcCtx.QuestionRepo.GetByIDs(l.ctx, questionIDs)
	if err != nil {
		l.Errorf("failed to get questions: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}
	qMap := make(map[string]*gormdb.Question)
	for _, q := range questions {
		qMap[q.ID] = q
	}

	score := 0
	total := len(req.Answers)
	wrongQuestionIDs := make([]string, 0, total)
	details := make([]types.SubmitExamDetailItem, 0, total)
	for qid, chosen := range req.Answers {
		opts := optsMap[qid]
		correct := false
		var correctOpt *gormdb.QuestionOption
		for _, o := range opts {
			if o.IsCorrect {
				correctOpt = o
				if o.OptionIndex == chosen {
					correct = true
					score++
				}
				break
			}
		}
		item := types.SubmitExamDetailItem{
			QuestionId:   qid,
			ChosenAnswer: chosen,
			Correct:      correct,
		}
		if !correct {
			wrongQuestionIDs = append(wrongQuestionIDs, qid)
		}
		if correctOpt != nil {
			item.CorrectOptionIndex = correctOpt.OptionIndex
			item.CorrectOptionDe = correctOpt.OptionDe
			item.CorrectOptionCn = correctOpt.OptionCn
		}
		if q := qMap[qid]; q != nil {
			item.Explanation = q.Explanation
		}
		details = append(details, item)
	}

	var stateID *string
	if req.StateId != "" {
		stateID = &req.StateId
	}

	record := &gormdb.ExamRecord{
		UserID:  userID,
		StateID: stateID,
		Total:   total,
		Score:   score,
		Passed:  score >= passScore,
		Answers: gormdb.ExamAnswers(req.Answers),
	}

	if err := l.svcCtx.ProgressRepo.CreateExamRecordWithWrongQuestions(l.ctx, record, wrongQuestionIDs); err != nil {
		l.Errorf("failed to create exam record: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	resp = &types.SubmitExamResponse{}
	resp.Base = *code.BaseSuccessResp()
	resp.Data.Id = record.ID
	resp.Data.Total = total
	resp.Data.Score = score
	resp.Data.Passed = record.Passed
	resp.Data.Details = details
	resp.Data.CreatedAt = record.CreatedAt.UnixMilli()
	return resp, nil
}
