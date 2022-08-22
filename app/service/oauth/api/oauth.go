package main

import (
	apollo "gohu/app/common/config"
	"gohu/app/common/log"
	"gohu/app/service/oauth/api/internal/config"
	"gohu/app/service/oauth/api/internal/handler"
	"gohu/app/service/oauth/api/internal/svc"
	"gohu/app/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

const serviceName = "oauth.api"

var c config.Config

func main() {
	// 初始化日志管理器
	_ = log.NewLogger()
	logWriter, _ := log.GetZapWriter()
	logx.MustSetup(log.GetLogXConfig(utils.GetServiceFullName(serviceName), "info"))
	logx.SetWriter(logWriter)

	logger := log.GetSugaredLogger()

	// 初始化配置管理器
	configClient, err := apollo.NewConfigClient()
	if err != nil {
		logger.Panicf("Initialize Apollo Client failed, err: %v", err)
	}

	// 初始化微服务设置
	namespace, serviceType, serviceSingleName := utils.GetServiceDetails(serviceName)
	err = configClient.UnmarshalServiceConfig(namespace, serviceType, serviceSingleName, &c)
	if err != nil {
		logger.Panicf("UnmarshalKey into service config failed, err: %v", err)
	}

	// 启动微服务服务器
	server := rest.MustNewServer(c.RestConf)

	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	logger.Infof("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
