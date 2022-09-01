package svc

import (
	"github.com/hibiken/asynq"
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/service/user/dao/query"
	"main/app/service/user/rpc/crud/internal/config"

	"github.com/go-redis/redis/v8"
	_ "github.com/spf13/viper/remote"
)

type ServiceContext struct {
	Config config.Config

	UserModel *query.Query
	Rdb       *redis.Client

	AsynqClient *asynq.Client

	ClientId     string
	ClientSecret string
}

const clientId = "default"

func NewServiceContext(c config.Config) *ServiceContext {
	logger := log.GetSugaredLogger()

	// 读取远程配置文件
	db, err := apollo.GetMysqlDB("user.yaml")
	if err != nil {
		logger.Fatalf("initialize mysql failed, err: %v", err)
	}

	rdb, err := apollo.GetRedisClient("user.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	clientSecret, err := apollo.GetClientSecret(clientId)
	if err != nil {
		logger.Fatalf("get client secret failed, err: %v", err)
	}
	return &ServiceContext{
		Config: c,

		UserModel: query.Use(db),
		Rdb:       rdb,

		AsynqClient: asynq.NewClient(c.AsynqClientConf),

		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}
