package svc

import (
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/service/oauth/rpc/token/store/internal/config"

	"github.com/go-redis/redis/v8"
)

type ServiceContext struct {
	Config config.Config

	Rdb *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	logger := log.GetSugaredLogger()

	rdb, err := apollo.GetRedisClient("oauth.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	return &ServiceContext{
		Config: c,
		Rdb:    rdb,
	}
}
