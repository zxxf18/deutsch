// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"

	"deutsch/internal/config"
	"deutsch/internal/middleware"
	"deutsch/internal/pkg/blacklist"
	"deutsch/internal/pkg/passwordcrypto"
	"deutsch/internal/sso"
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
	SSO                 *sso.Service
	SSOMiddleware       rest.Middleware
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

	issuer := c.OIDC.Issuer
	if issuer == "" {
		issuer = os.Getenv("DEUTSCH_OIDC_ISSUER")
	}
	clientID := c.OIDC.ClientID
	if clientID == "" {
		clientID = os.Getenv("DEUTSCH_OIDC_CLIENT_ID")
	}
	clientSecret := c.OIDC.ClientSecret
	if clientSecret == "" {
		clientSecret = os.Getenv("DEUTSCH_OIDC_CLIENT_SECRET")
	}
	redirectURL := c.OIDC.RedirectURL
	if redirectURL == "" {
		redirectURL = os.Getenv("DEUTSCH_OIDC_REDIRECT_URL")
	}
	sessionSecret := c.OIDC.SessionSecret
	if sessionSecret == "" {
		sessionSecret = os.Getenv("DEUTSCH_OIDC_SESSION_SECRET")
	}
	cookieName := c.OIDC.CookieName
	if cookieName == "" {
		cookieName = "deutsch_session"
	}
	auth := sso.New(sso.Config{Issuer: issuer, ClientID: clientID, ClientSecret: clientSecret, RedirectURL: redirectURL, SessionSecret: sessionSecret, CookieName: cookieName, AdminEmails: c.OIDC.AdminEmails})
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
		SSO:                 auth,
		SSOMiddleware:       middleware.NewSSOMiddleware(auth, repository.NewUserGormRepo(gormdb.DB)).Handle,
	}, nil
}
