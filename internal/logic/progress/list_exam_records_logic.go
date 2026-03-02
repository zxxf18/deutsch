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

type ListExamRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListExamRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListExamRecordsLogic {
	return &ListExamRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListExamRecordsLogic) ListExamRecords(req *types.Filter) (resp *types.ListExamRecordsResponse, err error) {
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

	list, total, err := l.svcCtx.ProgressRepo.GetExamRecordsByUser(l.ctx, userID, offset, pageSize)
	if err != nil {
		l.Errorf("failed to list exam records: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	items := make([]types.ExamRecordItem, 0, len(list))
	for _, r := range list {
		stateID := ""
		if r.StateID != nil {
			stateID = *r.StateID
		}
		items = append(items, types.ExamRecordItem{
			Id:        r.ID,
			StateId:   stateID,
			Total:     r.Total,
			Score:     r.Score,
			Passed:    r.Passed,
			CreatedAt: r.CreatedAt.UnixMilli(),
		})
	}

	resp = &types.ListExamRecordsResponse{}
	resp.Base = *code.BaseSuccessResp()
	resp.Data.Total = total
	resp.Data.Items = items
	return resp, nil
}
