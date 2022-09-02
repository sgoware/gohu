package config

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf
	Consul consul.Conf

	AsynqClientConf asynq.RedisClientOpt
}
