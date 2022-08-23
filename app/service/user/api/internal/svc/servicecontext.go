package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"main/app/service/user/api/internal/config"
	"main/app/service/user/api/internal/middleware"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware().Handle,
	}
}
