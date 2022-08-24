package svc

import (
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/rest"
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/common/middleware"
	"main/app/service/user/api/internal/config"
	"main/app/service/user/rpc/crud/crud"
	"main/app/service/user/rpc/vip/vip"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	Domain        string
	CrudRpcClient crud.Crud
	VipRpcClient  vip.Vip
	Cookie        *apollo.CookieConfig
	Rdb           *redis.Client

	AuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	logger := log.GetSugaredLogger()
	configClient, err := apollo.GetConfigClient()
	if err != nil {
		logger.Errorf("get configClient failed, err: %v", err)
	}
	rdb := redis.NewClient(configClient.NewRedisOptions("user.yaml"))

	cookieConfig, err := configClient.NewCookieConfig()
	if err != nil {
		logger.Errorf("get cookieConfig failed, err: %v", err)
	}
	domain, err := configClient.GetDomain()
	if err != nil {
		logger.Errorf("get domain failed, err: %v", err)
	}
	return &ServiceContext{
		Config:        c,
		Domain:        domain,
		CrudRpcClient: crud.NewCrud(zrpc.MustNewClient(c.CrudRpcClientConf)),
		VipRpcClient:  vip.NewVip(zrpc.MustNewClient(c.VipRpcClientConf)),
		Cookie:        cookieConfig,

		AuthMiddleware: middleware.NewAuthMiddleware(domain, cookieConfig, rdb).Handle,
	}
}
