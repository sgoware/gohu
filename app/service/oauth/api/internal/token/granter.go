package token

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gohu/app/common/config"
	"gohu/app/service/oauth/model"
)

const (
	GrantByAuth         = "authorization"
	GrantByRefreshToken = "refresh_token"

	EnhancerRpcClientConf = "Api.TokenEnhancerRpcClientConf.Target"
	StoreRpcClientConf    = "Api.TokenStoreRpcClientConf.Target"
)

var (
	tokenService model.TokenService
	tokenGranter model.TokenGranter
)

func InitTokenService() {
	configClient, _ := config.GetConfigClient()
	v, _ := configClient.GetViper("oauth.yaml")
	tokenService = model.NewRpcTokenService(zrpc.RpcClientConf{Target: v.GetString(EnhancerRpcClientConf)},
		zrpc.RpcClientConf{Target: v.GetString(StoreRpcClientConf)})
}

func InitTokenGranter() {
	tokenGranter = model.NewComposeTokenGranter(map[string]model.TokenGranter{
		"authorization": model.NewAuthorizationTokenGranter(GrantByAuth,
			model.GetClientDetails(),
			tokenService),
		"refresh_token": model.NewRefreshGranter(GrantByRefreshToken, tokenService),
	})
}

func GetTokenService() model.TokenService {
	return tokenService
}

func GetTokenGranter() model.TokenGranter {
	return tokenGranter
}
