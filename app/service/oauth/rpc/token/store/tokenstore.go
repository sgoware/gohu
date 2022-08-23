package main

import (
	"github.com/zeromicro/go-zero/core/logx"
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/utils"

	"main/app/service/oauth/rpc/token/store/internal/config"
	"main/app/service/oauth/rpc/token/store/internal/server"
	"main/app/service/oauth/rpc/token/store/internal/svc"
	"main/app/service/oauth/rpc/token/store/pb"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serviceName = "oauth.rpc.tokenStore"

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
		pb.RegisterTokenStoreServer(grpcServer, server.NewTokenStoreServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	// 注册服务到consul
	_ = consul.RegisterService(c.ListenOn, c.Consul)

	defer s.Stop()

	logger.Infof("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
