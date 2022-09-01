package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"main/app/service/question/mq/config"
	"main/app/service/question/rpc/crud/crud"
)

type ServiceContext struct {
	Config config.Config

	QuestionCrudRpcClient crud.Crud
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		QuestionCrudRpcClient: crud.NewCrud(zrpc.MustNewClient(c.QuestionCrudRpcClientConf)),
	}
}
