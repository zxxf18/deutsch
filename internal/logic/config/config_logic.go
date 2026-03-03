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

// 应用配置常量（可从 DB 或配置文件迁移）
// 300 通用题 + 16 州 × 10 州题 = 460
const (
	TotalQuestions     = 460
	ExamQuestions      = 33
	ExamMinutes        = 30
	PassScore          = 17
	TrialQuestionCount = 10
)

var languageModes = []types.LanguageModeItem{
	{Value: "de", Label: "Deutsch"},
	{Value: "de-cn", Label: "Deutsch + 中文"},
}

type ConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigLogic {
	return &ConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigLogic) Config() (resp *types.GetConfigResponse, err error) {
	data, err := configcache.GetAppConfig(l.svcCtx.Redis, func() ([]byte, error) {
		resp := &types.GetConfigResponse{}
		resp.Base = *code.BaseSuccessResp()
		resp.Data = types.AppConfigData{
			TotalQuestions:     TotalQuestions,
			ExamQuestions:      ExamQuestions,
			ExamMinutes:        ExamMinutes,
			PassScore:          PassScore,
			TrialQuestionCount: TrialQuestionCount,
			LanguageModes:      languageModes,
		}
		return json.Marshal(resp)
	})
	if err != nil {
		l.Errorf("failed to get config: %+v", err)
		return nil, code.NewCodeError(code.CodeInternalServerError)
	}
	var out types.GetConfigResponse
	if err := json.Unmarshal(data, &out); err != nil {
		l.Errorf("failed to unmarshal config cache: %+v", err)
		return nil, code.NewCodeError(code.CodeInternalServerError)
	}
	return &out, nil
}
