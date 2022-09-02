package svc

import (
	"github.com/hibiken/asynq"
	"main/app/common/log"
	"main/app/service/mq/asynq/processor/internal/config"
)

type ServiceContext struct {
	Config config.Config

	AsynqServer *asynq.Server
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: config.Config{},

		AsynqServer: newAsynqServer(c),
	}
}

func newAsynqServer(c config.Config) *asynq.Server {
	logger := log.GetSugaredLogger()
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: c.Redis.Host, Password: c.Redis.Pass},
		asynq.Config{
			Logger:      logger,
			LogLevel:    asynq.DebugLevel,
			Concurrency: 20, //max concurrent process job task num
		},
	)
}
