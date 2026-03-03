// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	JWTAuth struct {
		AccessSecret string
		AccessExpire int64 // 单位是s
	}
	MySQL struct {
		DataSource string // DSN: user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True
	}
	Redis     redis.RedisConf
	AssetsDir string `json:",optional"` // 静态资源目录，默认 assets
}
