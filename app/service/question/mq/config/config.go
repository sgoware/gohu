package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"main/app/common/mq/nsq"
)

type Config struct {
	service.ServiceConf

	NsqConsumerConf nsq.ConsumerConf

	QuestionCrudRpcClientConf zrpc.RpcClientConf
}
