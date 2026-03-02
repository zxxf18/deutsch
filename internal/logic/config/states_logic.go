package config

import (
	"context"
	"encoding/json"

	"deutsch/internal/code"
	"deutsch/internal/pkg/configcache"
	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatesLogic {
	return &StatesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatesLogic) States() (resp *types.ListStatesResponse, err error) {
	data, err := configcache.GetStates(l.svcCtx.Redis, func() ([]byte, error) {
		states, err := l.svcCtx.ConfigRepo.ListStates(l.ctx)
		if err != nil {
			return nil, err
		}
		resp := &types.ListStatesResponse{}
		resp.Base = *code.BaseSuccessResp()
		resp.Data.Total = int64(len(states))
		resp.Data.Items = make([]types.StateItem, 0, len(states))
		for _, s := range states {
			resp.Data.Items = append(resp.Data.Items, types.StateItem{
				Id:     s.ID,
				Slug:   s.Slug,
				Name:   s.Name,
				NameCn: s.NameCn,
			})
		}
		return json.Marshal(resp)
	})
	if err != nil {
		l.Errorf("failed to list states: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}
	var out types.ListStatesResponse
	if err := json.Unmarshal(data, &out); err != nil {
		l.Errorf("failed to unmarshal states cache: %+v", err)
		return nil, code.NewCodeError(code.CodeInternalServerError)
	}
	return &out, nil
}
