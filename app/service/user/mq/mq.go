package main

import (
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/service/user/mq/config"
	"main/app/service/user/mq/listen"
	"main/app/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

const serviceName = "user.mq"

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
	logx.MustSetup(log.GetLogXConfig(utils.GetServiceFullName(serviceName), "info"))
	logx.SetWriter(logWriter)

	// 初始化配置管理器
	configClient, err := apollo.NewConfigClient()
	if err != nil {
		logger.Panicf("initialize Apollo Client failed, err: %v", err)
	}

	// 初始化消息队列设置
	namespace, serviceType, serviceSingleName := utils.GetServiceDetails(serviceName)
	err = configClient.UnmarshalServiceConfig(namespace, serviceType, serviceSingleName, &c)
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
	mqs, err := listen.Mqs(c)
	if err != nil {
		logger.Fatalf("listen services failed, err: %v", err)
	}
	for _, mq := range mqs {
		serviceGroup.Add(mq)
	}

	serviceGroup.Start()
}
