package main

import (
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/service/oauth/api/internal/config"
	"main/app/service/oauth/api/internal/handler"
	"main/app/service/oauth/api/internal/svc"
	"main/app/service/oauth/api/internal/token"
	"main/app/service/oauth/model"
	"main/app/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

const serviceName = "oauth.api"

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
		logger.Fatalf("Initialize Apollo Client failed, err: %v", err)
	}

	err = model.InitClientDetails()
	if err != nil {
		logger.Fatalf("initialize client details failed, err: %v", err)
	}

	err = token.InitTokenService()
	if err != nil {
		logger.Fatalf("initialize token service failed, err: %v", err)
	}
	err = token.InitTokenGranter()
	if err != nil {
		logger.Fatalf("initialize token granter failed, err: %v", err)
	}

	// 初始化微服务设置
	namespace, serviceType, serviceSingleName := utils.GetServiceDetails(serviceName)
	err = configClient.UnmarshalServiceConfig(namespace, serviceType, serviceSingleName, &c)
	if err != nil {
		logger.Fatalf("UnmarshalKey into service config failed, err: %v", err)
	}

	// 启动微服务服务器
	server := rest.MustNewServer(c.RestConf)

	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	logger.Infof("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
