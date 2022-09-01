package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"main/app/service/notification/api/internal/config"
	"main/app/service/notification/rpc/info/info"
)

type ServiceContext struct {
	Config config.Config

	InfoRpcClient info.Info
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		InfoRpcClient: info.NewInfo(zrpc.MustNewClient(c.InfoRpcClientConf)),
	}
}
