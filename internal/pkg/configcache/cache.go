package configcache

import (
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	keyStates = "deutsch:config:states:v1"
	keyApp    = "deutsch:config:app:v1"
	ttl       = time.Hour // 1 小时，配置类数据变更少
)

// GetStates 先查缓存，未命中再调用 loadFromDB
func GetStates(rds *redis.Redis, loadFromDB func() ([]byte, error)) ([]byte, error) {
	val, err := rds.Get(keyStates)
	if err == nil && val != "" {
		return []byte(val), nil
	}
	data, err := loadFromDB()
	if err != nil {
		return nil, err
	}
	_ = rds.Setex(keyStates, string(data), int(ttl.Seconds()))
	return data, nil
}

// GetAppConfig 先查缓存，未命中再调用 build
func GetAppConfig(rds *redis.Redis, build func() ([]byte, error)) ([]byte, error) {
	val, err := rds.Get(keyApp)
	if err == nil && val != "" {
		return []byte(val), nil
	}
	data, err := build()
	if err != nil {
		return nil, err
	}
	_ = rds.Setex(keyApp, string(data), int(ttl.Seconds()))
	return data, nil
}
