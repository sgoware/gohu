package nsq

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	apollo "main/app/common/config"
	"main/app/common/mq/nsq"
	"main/app/service/notification/mq/internal/config"
	"main/app/service/notification/mq/internal/svc"
)

func NewService(c config.Config, svcContext *svc.ServiceContext) ([]service.Service, error) {
	domain, err := apollo.GetDomain()
	if err != nil {
		return nil, fmt.Errorf("get domain failed, %v", err)
	}

	nsqConsumerService, err := nsq.NewConsumerService(
		c.NsqConsumerConf.Topic,
		c.NsqConsumerConf.Channel,
		&PublishNotificationHandler{
			Domain:        domain,
			CrudRpcClient: svcContext.CrudRpcClient,
		},
	)
	if err != nil {
		return nil, err
	}

	return []service.Service{
		nsqConsumerService,
	}, nil
}
