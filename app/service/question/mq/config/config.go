package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"main/app/common/mq/nsq"
)

type Config struct {
	rest.RestConf

	NsqConsumerConf nsq.ConsumerConf

	QuestionCrudRpcClientConf zrpc.RpcClientConf
}
