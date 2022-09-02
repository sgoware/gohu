package config

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/service"
)

type Config struct {
	service.ServiceConf

	RedisConf asynq.RedisClientOpt
}
