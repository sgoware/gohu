package svc

import (
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/service/user/dao/query"
	"main/app/service/user/rpc/vip/internal/config"
)

type ServiceContext struct {
	Config config.Config

	UserModel *query.Query
}

func NewServiceContext(c config.Config) *ServiceContext {
	logger := log.GetSugaredLogger()

	db, err := apollo.GetMysqlDB("user.yaml")
	if err != nil {
		logger.Fatalf("initialize mysql failed, err: %v", err)
	}

	return &ServiceContext{
		Config:    c,
		UserModel: query.Use(db),
	}
}
