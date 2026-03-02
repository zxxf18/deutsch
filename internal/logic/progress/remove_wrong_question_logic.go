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
)

type RemoveWrongQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveWrongQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveWrongQuestionLogic {
	return &RemoveWrongQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveWrongQuestionLogic) RemoveWrongQuestion(req *types.RemoveWrongQuestionRequest) (resp *types.RemoveWrongQuestionResponse, err error) {
	userID := common.GetUserID(l.ctx)
	if userID == "" {
		return nil, code.NewCodeError(code.CodeUnauthorized)
	}

	if req.QuestionId == "" {
		return nil, code.NewCodeErrorWithMsg(code.CodeValidationError, "题目ID不能为空")
	}

	if err := l.svcCtx.ProgressRepo.RemoveWrongQuestion(l.ctx, userID, req.QuestionId); err != nil {
		l.Errorf("failed to remove wrong question: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	resp = &types.RemoveWrongQuestionResponse{}
	resp.Base = *code.BaseSuccessResp()
	return resp, nil
}
