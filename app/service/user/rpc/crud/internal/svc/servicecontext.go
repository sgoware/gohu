package svc

import (
	"context"
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/service/user/dao/query"
	"main/app/service/user/rpc/crud/internal/config"

	"github.com/go-redis/redis/v8"
	_ "github.com/spf13/viper/remote"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	UserModel *query.Query
	Rdb       *redis.Client

	ClientId     string
	ClientSecret string
}

func NewServiceContext(c config.Config) *ServiceContext {
	logger := log.GetSugaredLogger()

	// 读取远程配置文件
	configClient, err := apollo.GetConfigClient()
	if err != nil {
		logger.Fatalf("get configClient failed, err: %v", err)
	}

	dsn, err := configClient.GetMysqlDsn("user.yaml")
	if err != nil {
		logger.Fatalf("get mysql dsn failed, err: %v", err)
	}

	// 连接mysql和redis
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		logger.Fatalf("initiate mysql failed, err: %v", err)
	}

	rdb := redis.NewClient(configClient.NewRedisOptions("user.yaml"))
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatalf("initiate redis failed, err: %v", err)
	}

	v, err := configClient.GetViper("oauth.yaml")
	if err != nil {
		logger.Fatalf("get viper failed, err: %v", err)
	}
	clientId := "default"
	clientSecret := v.GetString("Client.default.Secret")
	return &ServiceContext{
		Config:       c,
		UserModel:    query.Use(db),
		Rdb:          rdb,
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}
