// Package jwt jwt插件
package jwt

import (
	"context"
	"time"
)

type JWTInfo struct {
	Token      string    `json:"token"`
	Expire     time.Time `json:"expire"`
	MaxRefresh time.Time `json:"maxRefresh"`
}

type JWT interface {
	GetJWT(c context.Context) (string, error)
	GenerateJWT(c context.Context) (*JWTInfo, error)
	RefreshJWT(c context.Context) (*JWTInfo, error)
	CheckAndParseJWT(c context.Context) (map[string]interface{}, error)
}
