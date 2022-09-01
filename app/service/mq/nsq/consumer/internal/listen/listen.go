package listen

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"main/app/service/mq/nsq/consumer/internal/config"
	"main/app/service/mq/nsq/consumer/internal/listen/user"
	"main/app/service/mq/nsq/consumer/internal/svc"
)

func NewServices(c config.Config) ([]service.Service, error) {
	var services []service.Service

	svcContext := svc.NewServiceContext(c)

	userServices, err := user.NewService(c, svcContext)
	if err != nil {
		return nil, fmt.Errorf("initialze nsq services failed, err: %v", err)
	}

	services = append(services, userServices...)

	return services, nil
}
