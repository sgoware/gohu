package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"main/app/service/mq/nsq/consumer/internal/config"
	notification "main/app/service/notification/rpc/crud/crud"
	user "main/app/service/user/rpc/crud/crud"
)

type ServiceContext struct {
	Config config.Config

	UserCrudRpcClient         user.Crud
	NotificationCrudRpcClient notification.Crud
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		UserCrudRpcClient: user.NewCrud(zrpc.MustNewClient(c.UserCrudRpcClientConf)),
	}
}
