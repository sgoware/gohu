package user

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	apollo "main/app/common/config"
	"main/app/common/mq/nsq"
	"main/app/service/mq/nsq/consumer/internal/config"
	"main/app/service/mq/nsq/consumer/internal/svc"
)

func NewService(c config.Config, svcContext *svc.ServiceContext) ([]service.Service, error) {
	domain, err := apollo.GetDomain()
	if err != nil {
		return nil, fmt.Errorf("get domain failed, %v", err)
	}
	publishNotificationConsumerService, err := nsq.NewConsumerService(
		c.PublishNotificationConsumerConf.Topic,
		c.PublishNotificationConsumerConf.Channel,
		&PublishNotificationHandler{
			Domain:        domain,
			CrudRpcClient: svcContext.NotificationCrudRpcClient,
		},
	)

	if err != nil {
		return nil, err
	}

	return []service.Service{
		publishNotificationConsumerService,
	}, nil
}
