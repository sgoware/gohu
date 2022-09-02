package svc

import (
	"github.com/go-redis/redis/v8"
	"github.com/yitter/idgenerator-go/idgen"
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/service/notification/dao/query"
	"main/app/service/notification/rpc/crud/internal/config"
)

type ServiceContext struct {
	Config            config.Config
	NotificationModel *query.Query
	Rdb               *redis.Client

	IdGenerator *idgen.DefaultIdGenerator
}

func NewServiceContext(c config.Config) *ServiceContext {
	logger := log.GetSugaredLogger()

	db, err := apollo.GetMysqlDB("notification.yaml")
	if err != nil {
		logger.Fatalf("initialize mysql failed, err: %v", err)
	}

	rdb, err := apollo.GetRedisClient("notification.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	idGenerator, err := apollo.NewIdGenerator("notification.yaml")
	if err != nil {
		logger.Fatalf("initialize idGenerator failed, err: %v", err)
	}

	return &ServiceContext{
		Config: c,

		NotificationModel: query.Use(db),
		Rdb:               rdb,

		IdGenerator: idGenerator,
	}
}
