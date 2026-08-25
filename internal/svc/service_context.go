// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"path/filepath"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"

	"deutsch/internal/config"
	"deutsch/internal/middleware"
	"deutsch/internal/pkg/blacklist"
	"deutsch/internal/pkg/passwordcrypto"
	"deutsch/model/gormdb"
	"deutsch/model/repository"
)

type ServiceContext struct {
	Config              config.Config
	AssetsDir           string
	Redis               *redis.Redis
	TokenBlacklist      *blacklist.TokenBlacklist
	JWTMiddleware       rest.Middleware
	BlacklistMiddleware rest.Middleware
	AdminMiddleware     rest.Middleware
	UserRepo            repository.UserRepository
	InviteCodeRepo      repository.InviteCodeRepository
	ConfigRepo          repository.ConfigRepository
	QuestionRepo        repository.QuestionRepository
	ProgressRepo        repository.ProgressRepository
	PasswordCipher      *passwordcrypto.Cipher
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	rds := redis.MustNewRedis(c.Redis)
	tokenBlacklist := blacklist.NewTokenBlacklist(rds)
	passwordCipher, err := passwordcrypto.New(c.PasswordEncryption.Key)
	if err != nil {
		return nil, err
	}

	assetsDir := c.AssetsDir
	if assetsDir == "" {
		assetsDir = "assets"
	}
	assetsDir, _ = filepath.Abs(assetsDir)

	return &ServiceContext{
		Config:              c,
		AssetsDir:           assetsDir,
		Redis:               rds,
		TokenBlacklist:      tokenBlacklist,
		JWTMiddleware:       middleware.NewJWTMiddleware(c).Handle,
		BlacklistMiddleware: middleware.NewBlacklistMiddleware(tokenBlacklist).Handle,
		AdminMiddleware:     middleware.NewAdminMiddleware(c).Handle,
		UserRepo:            repository.NewUserGormRepo(gormdb.DB),
		InviteCodeRepo:      repository.NewInviteGormRepo(gormdb.DB),
		ConfigRepo:          repository.NewConfigGormRepo(gormdb.DB),
		QuestionRepo:        repository.NewQuestionGormRepo(gormdb.DB),
		ProgressRepo:        repository.NewProgressGormRepo(gormdb.DB),
		PasswordCipher:      passwordCipher,
	}, nil
}
