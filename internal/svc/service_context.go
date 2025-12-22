// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/zeromicro/go-zero/rest"

	"deutsch/internal/config"
	"deutsch/internal/middleware"
)

type ServiceContext struct {
	Config          config.Config
	JWTMiddleware   rest.Middleware
	AdminMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:          c,
		JWTMiddleware:   middleware.NewJWTMiddleware(c).Handle,
		AdminMiddleware: middleware.NewAdminMiddleware().Handle,
	}
}
