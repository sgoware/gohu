package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	service.ServiceConf

	Redis redis.RedisConf

	UserCrudRpcClientConf     zrpc.RpcClientConf
	QuestionCrudRpcClientConf zrpc.RpcClientConf
	CommentCrudRpcClientConf  zrpc.RpcClientConf
}
