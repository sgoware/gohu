package svc

import (
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/service/oauth/rpc/token/enhancer/internal/config"
	"main/app/service/oauth/rpc/token/enhancer/internal/jwt"
	"main/app/service/oauth/rpc/token/store/tokenstore"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	TokenStoreRpcClient tokenstore.TokenStore

	Enhancer *jwt.JWT
}

func NewServiceContext(c config.Config) *ServiceContext {
	logger := log.GetSugaredLogger()

	configClient, err := apollo.GetConfigClient()
	if err != nil {
		logger.Errorf("get configClient failed, err: %v", err)
	}
	v, err := configClient.GetViper("oauth.yaml")
	if err != nil {
		logger.Errorf("get viper failed, err: %v", err)
	}
	return &ServiceContext{
		Config:              c,
		TokenStoreRpcClient: tokenstore.NewTokenStore(zrpc.MustNewClient(c.TokenStoreRpcClientConf)),
		Enhancer: jwt.NewJWT(v.GetString("JWTAuth.Secret"),
			v.GetString("JWTAuth.Issuer")),
	}
}
