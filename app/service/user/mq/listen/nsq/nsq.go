package nsq

import (
	"github.com/zeromicro/go-zero/core/service"
	"main/app/common/mq/nsq"
	"main/app/service/user/mq/config"
	"main/app/service/user/mq/svc"
)

func NewService(c config.Config, svcContext *svc.ServiceContext) ([]service.Service, error) {
	nsqConsumerService, err := nsq.NewConsumerService(
		c.NsqConsumerConf.Topic,
		c.NsqConsumerConf.Channel,
		&ChangeFollowerHandler{
			CrudRpcClient: svcContext.UserCrudRpcClient,
		},
	)
	if err != nil {
		return nil, err
	}

	return []service.Service{
		nsqConsumerService,
	}, nil
}
