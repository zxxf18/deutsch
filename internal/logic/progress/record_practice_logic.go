// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package progress

import (
	"context"

	"deutsch/internal/code"
	"deutsch/internal/common"
	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RecordPracticeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecordPracticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecordPracticeLogic {
	return &RecordPracticeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecordPracticeLogic) RecordPractice(req *types.RecordPracticeRequest) (resp *types.RecordPracticeResponse, err error) {
	userID := common.GetUserID(l.ctx)
	if userID == "" {
		return nil, code.NewCodeError(code.CodeUnauthorized)
	}

	if req.QuestionId == "" {
		return nil, code.NewCodeErrorWithMsg(code.CodeValidationError, "题目ID不能为空")
	}

	_, err = l.svcCtx.QuestionRepo.GetByID(l.ctx, req.QuestionId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.NewCodeError(code.CodeQuestionNotFound)
		}
		l.Errorf("failed to get question: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	if err := l.svcCtx.ProgressRepo.RecordPractice(l.ctx, userID, req.QuestionId, req.Correct); err != nil {
		l.Errorf("failed to upsert progress: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	resp = &types.RecordPracticeResponse{}
	resp.Base = *code.BaseSuccessResp()
	return resp, nil
}
