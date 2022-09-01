package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"main/app/service/user/mq/config"
	"main/app/service/user/rpc/crud/crud"
)

type ServiceContext struct {
	Config config.Config

	UserCrudRpcClient crud.Crud
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		UserCrudRpcClient: crud.NewCrud(zrpc.MustNewClient(c.UserCrudRpcClientConf)),
	}
}
