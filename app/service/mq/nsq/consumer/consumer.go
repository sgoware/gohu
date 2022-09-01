package consumer

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/service/mq/nsq/consumer/internal/config"
	"main/app/service/mq/nsq/consumer/internal/listen"
)

const mqName = "mq.nsq.consumer"

var c config.Config

func main() {
	// 初始化日志管理器
	err := log.InitLogger()
	if err != nil {
		panic("initialize logger failed")
	}
	logger := log.GetSugaredLogger()

	logWriter, err := log.GetZapWriter()
	if err != nil {
		logger.Fatalf("get log writer failed")
	}
	logx.MustSetup(log.GetLogXConfig("mq-nsq-consumer", "info"))
	logx.SetWriter(logWriter)

	// 初始化配置管理器
	configClient, err := apollo.NewConfigClient()
	if err != nil {
		logger.Panicf("initialize Apollo Client failed, err: %v", err)
	}

	// 初始化消息队列设置
	err = configClient.UnmarshalKey("mq", "nsq.consumer", &c)
	if err != nil {
		logger.Fatalf("UnmarshalKey into service config failed, err: %v", err)
	}

	// 初始化log、trace
	err = c.SetUp()
	if err != nil {
		logger.Fatalf("initialize go-zero internal service failed, err: %v")
	}

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	consumerServices, err := listen.NewServices(c)
	if err != nil {
		logger.Fatalf("initialize nsq consumer services failed, err: %v")
	}

	for _, consumerService := range consumerServices {
		serviceGroup.Add(consumerService)
	}

	serviceGroup.Start()
}
