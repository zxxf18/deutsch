package jwt

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4/request"
	"github.com/zeromicro/go-zero/core/logx"
)

type JWTConfig struct {
	SigningAlgorithm string        `yaml:"sa" json:"sa" default:"HS256"`
	Key              string        `yaml:"key" json:"key" default:"Deutsch.20251014"`
	PrivKeyFile      string        `yaml:"privKeyFile" json:"privKeyFile"`
	PubKeyFile       string        `yaml:"pubKeyFile" json:"pubKeyFile"`
	Timeout          time.Duration `yaml:"timeout" json:"timeout" default:"30m"`
	MaxRefresh       time.Duration `yaml:"maxRefresh" json:"maxRefresh" default:"1h"`
}

type defaultJWT struct {
	ctx    context.Context
	cfg    JWTConfig
	helper *JWTHelper
}

func DefaultJWTConfig(secret string, expire int64) *JWTConfig {
	return &JWTConfig{
		Key:        secret,
		Timeout:    time.Duration(expire) * time.Second,
		MaxRefresh: time.Duration(expire*2) * time.Second,
	}
}

func NewWithConfig(ctx context.Context, cfg JWTConfig) (any, error) {
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

func New(ctx context.Context, secret string, expire int64) (any, error) {
	cfg := DefaultJWTConfig(secret, expire)
	helper, err := NewJWTHelper(*cfg)
	if err != nil {
		return nil, err
	}
	return &defaultJWT{
		cfg:    *cfg,
		helper: helper,
		ctx:    logx.ContextWithFields(ctx, logx.Field("pkg", "defaultJWT")),
	}, nil
}

func (d *defaultJWT) GenerateJWT(data any) (*JWTInfo, error) {
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

func (d *defaultJWT) RefreshJWT(r *http.Request) (*JWTInfo, error) {
	j, exp, err := d.helper.Refresh(r)
	if err != nil {
		return nil, err
	}
	return &JWTInfo{
		Token:      j,
		Expire:     exp,
		MaxRefresh: exp.Add(d.helper.MaxRefresh - d.helper.Timeout),
	}, nil
}

// GetJWT
// 采用和 go-zero 框架的相同的存储和解析位置逻辑
// 为从 header 的 Authorization 中获取 token 并解析
func (d *defaultJWT) GetJWT(r *http.Request) (string, error) {
	return request.AuthorizationHeaderExtractor.ExtractToken(r)
}

func (d *defaultJWT) CheckAndParseJWT(r *http.Request) (map[string]any, error) {
	return d.helper.CheckExpireAndParse(r)
}
