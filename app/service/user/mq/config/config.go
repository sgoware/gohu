package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"main/app/common/mq/nsq"
)

type Config struct {
	rest.RestConf // rest api配置

	NsqConsumerConf nsq.ConsumerConf

	// rpc client配置
	UserCrudRpcClientConf zrpc.RpcClientConf
}
