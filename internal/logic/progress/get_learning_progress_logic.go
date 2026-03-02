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

type GetLearningProgressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLearningProgressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLearningProgressLogic {
	return &GetLearningProgressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLearningProgressLogic) GetLearningProgress() (resp *types.GetLearningProgressResponse, err error) {
	userID := common.GetUserID(l.ctx)
	if userID == "" {
		return nil, code.NewCodeError(code.CodeUnauthorized)
	}

	list, err := l.svcCtx.ProgressRepo.GetProgressByUser(l.ctx, userID)
	if err != nil {
		l.Errorf("failed to get progress: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	summaryMap := make(map[string]*types.LearningProgressSummary)
	for _, p := range list {
		q, err := l.svcCtx.QuestionRepo.GetByID(l.ctx, p.QuestionID)
		if err != nil || q == nil {
			continue
		}
		stateID := gormdb.GeneralStateID
		if q.StateID != nil && *q.StateID != "" {
			stateID = *q.StateID
		}
		if s, ok := summaryMap[stateID]; ok {
			s.PracticedCount++
			s.CorrectCount += p.CorrectCount
		} else {
			summaryMap[stateID] = &types.LearningProgressSummary{
				StateId:        stateID,
				PracticedCount: 1,
				CorrectCount:   p.CorrectCount,
			}
		}
	}

	// 填充每类题目总数
	for stateID, s := range summaryMap {
		total, err := l.svcCtx.QuestionRepo.CountByStateID(l.ctx, stateID)
		if err != nil {
			l.Errorf("failed to count questions for state %s: %+v", stateID, err)
			continue
		}
		s.Total = int(total)
	}

	items := make([]types.LearningProgressSummary, 0, len(summaryMap))
	for _, v := range summaryMap {
		items = append(items, *v)
	}

	resp = &types.GetLearningProgressResponse{}
	resp.Base = *code.BaseSuccessResp()
	resp.Data.Items = items
	return resp, nil
}
