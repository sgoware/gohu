package listen

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"main/app/service/comment/mq/internal/config"
	"main/app/service/comment/mq/internal/listen/nsq"
	"main/app/service/comment/mq/internal/svc"
)

func Mqs(c config.Config) ([]service.Service, error) {
	var services []service.Service

	svcContext := svc.NewServiceContext(c)

	nsqServices, err := nsq.NewService(c, svcContext)
	if err != nil {
		return nil, fmt.Errorf("initialze nsq services failed, err: %v", err)
	}

	services = append(services, nsqServices...)

	return services, nil
}
