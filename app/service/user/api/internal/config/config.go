package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf // rest api配置

	// rpc client配置
	CrudRpcClientConf zrpc.RpcClientConf
	VipRpcClientConf  zrpc.RpcClientConf
}
