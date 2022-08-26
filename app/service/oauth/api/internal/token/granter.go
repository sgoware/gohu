package token

import (
	"fmt"
	"github.com/zeromicro/go-zero/zrpc"
	"main/app/common/config"
	"main/app/service/oauth/model"
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

func InitTokenService() (err error) {
	configClient, err := config.GetConfigClient()
	if err != nil {
		return fmt.Errorf("get config client failed, %v", err)
	}

	v, err := configClient.GetViper("oauth.yaml")
	if err != nil {
		return fmt.Errorf("get viper failed, %v", err)
	}

	tokenService = model.NewRpcTokenService(zrpc.RpcClientConf{Target: v.GetString(EnhancerRpcClientConf)},
		zrpc.RpcClientConf{Target: v.GetString(StoreRpcClientConf)})

	return nil
}

func InitTokenGranter() (err error) {
	authorizationTokenGranter, err := model.NewAuthorizationTokenGranter(GrantByAuth,
		model.GetClientDetails(),
		tokenService)
	if err != nil {
		return err
	}

	refreshGranter, err := model.NewRefreshGranter(GrantByRefreshToken, tokenService)
	if err != nil {
		return err
	}

	tokenGranter = model.NewComposeTokenGranter(map[string]model.TokenGranter{
		"authorization": authorizationTokenGranter,
		"refresh_token": refreshGranter,
	})
	return nil
}

func GetTokenService() model.TokenService {
	return tokenService
}

func GetTokenGranter() model.TokenGranter {
	return tokenGranter
}
