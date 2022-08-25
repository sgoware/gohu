package svc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	configClient, err := apollo.GetConfigClient()
	if err != nil {
		logger.Fatalf("get configClient failed, err: %v", err)
	}

	dsn, err := configClient.GetMysqlDsn("user.yaml")
	if err != nil {
		logger.Fatalf("get mysql dsn failed, err: %v", err)
	}
	db, _ := gorm.Open(mysql.Open(dsn))
	return &ServiceContext{
		Config:    c,
		UserModel: query.Use(db),
	}
}
