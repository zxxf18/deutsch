// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"

	"deutsch/internal/config"
	"deutsch/internal/middleware"
	"deutsch/internal/pkg/blacklist"
	"deutsch/model/gormdb"
	"deutsch/model/repository"
)

type ServiceContext struct {
	Config              config.Config
	Redis               *redis.Redis
	TokenBlacklist      *blacklist.TokenBlacklist
	JWTMiddleware       rest.Middleware
	BlacklistMiddleware rest.Middleware
	AdminMiddleware     rest.Middleware
	UserRepo            repository.UserRepository
	InviteCodeRepo      repository.InviteCodeRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds := redis.MustNewRedis(c.Redis)
	tokenBlacklist := blacklist.NewTokenBlacklist(rds)

	return &ServiceContext{
		Config:              c,
		Redis:               rds,
		TokenBlacklist:      tokenBlacklist,
		JWTMiddleware:       middleware.NewJWTMiddleware(c).Handle,
		BlacklistMiddleware: middleware.NewBlacklistMiddleware(tokenBlacklist).Handle,
		AdminMiddleware:     middleware.NewAdminMiddleware(c).Handle,
		UserRepo:            repository.NewUserGormRepo(gormdb.DB),
		InviteCodeRepo:      repository.NewInviteGormRepo(gormdb.DB),
	}
}
