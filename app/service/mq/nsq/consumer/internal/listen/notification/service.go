package notification

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	apollo "main/app/common/config"
	"main/app/common/mq/nsq"
	"main/app/service/mq/nsq/consumer/internal/config"
	"main/app/service/mq/nsq/consumer/internal/svc"
	"main/app/service/notification/rpc/crud/crud"
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
			Domain:                    domain,
			NotificationCrudRpcClient: crud.NewCrud(zrpc.MustNewClient(c.NotificationCrudRpcClientConf)),
		},
	)

	if err != nil {
		return nil, err
	}

	return []service.Service{
		publishNotificationConsumerService,
	}, nil
}
