package main

import (
	"github.com/zeromicro/go-zero/core/logx"
	apollo "gohu/app/common/config"
	"gohu/app/common/log"
	"gohu/app/utils"

	"gohu/app/service/oauth/rpc/token/enhancer/internal/config"
	"gohu/app/service/oauth/rpc/token/enhancer/internal/server"
	"gohu/app/service/oauth/rpc/token/enhancer/internal/svc"
	"gohu/app/service/oauth/rpc/token/enhancer/pb"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serviceName = "oauth.rpc.tokenEnhancer"

var c config.Config

func main() {
	// 初始化日志管理器
	_ = log.NewLogger()
	lowWriter, _ := log.GetZapWriter()
	logx.MustSetup(log.GetLogXConfig(utils.GetServiceFullName(serviceName), "info"))
	logx.SetWriter(lowWriter)

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

	ctx := svc.NewServiceContext(c)

	// 启动微服务服务器
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterTokenEnhancerServer(grpcServer, server.NewTokenEnhancerServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	logger.Infof("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
