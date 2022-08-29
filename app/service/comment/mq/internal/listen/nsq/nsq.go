package nsq

import (
	"github.com/zeromicro/go-zero/core/service"
	"main/app/common/mq/nsq"
	"main/app/service/comment/mq/internal/config"
	"main/app/service/comment/mq/internal/svc"
)

func NewService(c config.Config, svcContext *svc.ServiceContext) ([]service.Service, error) {
	nsqConsumerService, err := nsq.NewConsumerService(
		c.NsqConsumerConf.Topic,
		c.NsqConsumerConf.Channel,
		&CommentSubjectHandler{
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
