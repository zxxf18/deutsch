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
	"gorm.io/gorm"
)

type GetExamRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetExamRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExamRecordLogic {
	return &GetExamRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetExamRecordLogic) GetExamRecord(req *types.GetExamRecordRequest) (resp *types.GetExamRecordResponse, err error) {
	userID := common.GetUserID(l.ctx)
	if userID == "" {
		return nil, code.NewCodeError(code.CodeUnauthorized)
	}

	record, err := l.svcCtx.ProgressRepo.GetExamRecordByID(l.ctx, req.Id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.NewCodeError(code.CodeExamRecordNotFound)
		}
		l.Errorf("failed to get exam record: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	questionIDs := make([]string, 0, len(record.Answers))
	for qid := range record.Answers {
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

	details := make([]types.ExamRecordDetailItem, 0, len(record.Answers))
	for qid, chosen := range record.Answers {
		opts := optsMap[qid]
		correct := false
		var correctOpt *gormdb.QuestionOption
		for _, o := range opts {
			if o.IsCorrect {
				correctOpt = o
				if o.OptionIndex == chosen {
					correct = true
				}
				break
			}
		}
		item := types.ExamRecordDetailItem{
			QuestionId:   qid,
			ChosenAnswer: chosen,
			Correct:      correct,
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

	stateID := ""
	if record.StateID != nil {
		stateID = *record.StateID
	}

	resp = &types.GetExamRecordResponse{}
	resp.Base = *code.BaseSuccessResp()
	resp.Data.Id = record.ID
	resp.Data.StateId = stateID
	resp.Data.Total = record.Total
	resp.Data.Score = record.Score
	resp.Data.Passed = record.Passed
	resp.Data.Details = details
	resp.Data.CreatedAt = record.CreatedAt.UnixMilli()
	return resp, nil
}
