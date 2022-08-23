package svc

import (
	"context"
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
	configClient, err := apollo.GetConfigClient()
	if err != nil {
		logger.Errorf("Get configClient failed, err: %v", err)
	}

	rdb := redis.NewClient(configClient.NewRedisOptions("oauth.yaml"))
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		logger.Errorf("Initiate redis failed, err: %v", err)
	}
	return &ServiceContext{
		Config: c,
		Rdb:    rdb,
	}
}
