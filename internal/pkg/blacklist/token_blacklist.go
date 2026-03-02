package blacklist

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

const keyPrefix = "deutsch:jwt:blacklist:"

// TokenBlacklist JWT 黑名单，登出后 token 不可用
type TokenBlacklist struct {
	rds *redis.Redis
}

// NewTokenBlacklist 创建 TokenBlacklist
func NewTokenBlacklist(rds *redis.Redis) *TokenBlacklist {
	return &TokenBlacklist{rds: rds}
}

// AddWithExpire 将 token 加入黑名单，expireAt 为 token 过期时间戳（秒）
func (b *TokenBlacklist) AddWithExpire(token string, expireAt int64) error {
	if token == "" {
		return nil
	}
	ttl := time.Unix(expireAt, 0).Sub(time.Now())
	if ttl <= 0 {
		return nil // 已过期，无需加入
	}
	key := b.key(token)
	return b.rds.Setex(key, "1", int(ttl.Seconds()))
}

// IsBlacklisted 检查 token 是否在黑名单中
// 注意：go-zero Redis.Get 在 key 不存在时返回 ("", nil)，不是 redis.Nil
func (b *TokenBlacklist) IsBlacklisted(token string) (bool, error) {
	if token == "" {
		return false, nil
	}
	key := b.key(token)
	val, err := b.rds.Get(key)
	if err != nil {
		return false, err
	}
	// key 不存在时 go-zero 返回 val="" 且 err=nil，需据此判断
	return val != "", nil
}

func (b *TokenBlacklist) key(token string) string {
	h := sha256.Sum256([]byte(token))
	return keyPrefix + hex.EncodeToString(h[:])
}
