package jwt

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"deutsch/internal/pkg/utils"
)

type JWTConfig struct {
	SigningAlgorithm string        `yaml:"sa" json:"sa" default:"HS256"`
	Key              string        `yaml:"key" json:"key" default:"Deutsch.20251014"`
	PrivKeyFile      string        `yaml:"privKeyFile" json:"privKeyFile"`
	PubKeyFile       string        `yaml:"pubKeyFile" json:"pubKeyFile"`
	Timeout          time.Duration `yaml:"timeout" json:"timeout" default:"30m"`
	MaxRefresh       time.Duration `yaml:"maxRefresh" json:"maxRefresh" default:"1h"`
	JWTKey           string        `yaml:"jwtkey" json:"jwtkey" default:"jwt"`
}

type defaultJWT struct {
	ctx    context.Context
	cfg    JWTConfig
	helper *JWTHelper
}

func New(ctx context.Context) (any, error) {
	var cfg JWTConfig
	if err := utils.LoadConfig(&cfg); err != nil {
		return nil, err
	}
	helper, err := NewJWTHelper(cfg)
	if err != nil {
		return nil, err
	}
	return &defaultJWT{
		cfg:    cfg,
		helper: helper,
		ctx:    logx.ContextWithFields(ctx, logx.Field("pkg", "defaultJWT")),
	}, nil
}

func (d *defaultJWT) GenerateJWT(c context.Context, data any) (*JWTInfo, error) {
	claims := map[string]any{
		JWTKey: data,
	}
	j, exp, err := d.helper.Generate(claims)
	if err != nil {
		return nil, err
	}
	return &JWTInfo{
		Token:      j,
		Expire:     exp,
		MaxRefresh: exp.Add(d.helper.MaxRefresh - d.helper.Timeout),
	}, nil
}

// GenerateDefaultJWT
// 通过获取 context 中的固定的key来生成 jwt token
func (d *defaultJWT) GenerateDefaultJWT(c context.Context) (*JWTInfo, error) {
	info := c.GetUserInfo()
	return d.GenerateJWT(c, info)
}

func (d *defaultJWT) RefreshJWT(c context.Context) (*JWTInfo, error) {
	j, exp, err := d.helper.Refresh(c.Context)
	if err != nil {
		return nil, err
	}
	return &JWTInfo{
		Token:      j,
		Expire:     exp,
		MaxRefresh: exp.Add(d.helper.MaxRefresh - d.helper.Timeout),
	}, nil
}

func (d *defaultJWT) GetJWT(c context.Context) (string, error) {
	return d.helper.GetTokenString(c.Context)
}

func (d *defaultJWT) CheckAndParseJWT(c context.Context) (map[string]any, error) {
	return d.helper.CheckExpireAndParse(c.Context)
}
