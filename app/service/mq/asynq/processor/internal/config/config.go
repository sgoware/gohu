package config

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	service.ServiceConf

	RedisConf asynq.RedisClientOpt

	UserCrudRpcClientConf     zrpc.RpcClientConf
	QuestionCrudRpcClientConf zrpc.RpcClientConf
	CommentCrudRpcClientConf  zrpc.RpcClientConf
}
