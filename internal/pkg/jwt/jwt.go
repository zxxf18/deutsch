// Package jwt jwt插件
package jwt

import (
	"net/http"
	"time"
)

type JWTInfo struct {
	Token      string    `json:"token"`
	Expire     time.Time `json:"expire"`
	MaxRefresh time.Time `json:"maxRefresh"`
}

// JWT
// 因为 Auth 层采用 go-zero 框架的 github.com/zeromicro/go-zero/rest 包下的解析逻辑
// 所以这里也采用同样版本的 jwt 生成和校验逻辑
type JWT interface {
	GetJWT(r *http.Request) (string, error)
	GenerateJWT(data any) (*JWTInfo, error)
	RefreshJWT(r *http.Request) (*JWTInfo, error)
	CheckAndParseJWT(r *http.Request) (map[string]any, error)
}
