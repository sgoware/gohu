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
	Jwt       *apollo.JWTConfig
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
		logger.Errorf("get configClient failed, err: %v", err)
	}

	// 连接mysql和redis
	db, _ := gorm.Open(mysql.Open(configClient.GetMysqlDsn("user.yaml")))
	rdb := redis.NewClient(configClient.NewRedisOptions("user.yaml"))
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		logger.Errorf("initiate redis failed, err: %v", err)
	}

	v, err := configClient.GetViper("oauth.yaml")
	if err != nil {
		logger.Errorf("get viper failed, err: %v", err)
	}
	clientId := "default"
	clientSecret := v.GetString("Client.default.Secret")
	return &ServiceContext{
		Config: c,
		Jwt: &apollo.JWTConfig{
			SecretKey:   v.GetString("JWTAuth.Secret"),
			ExpiresTime: v.GetInt64("JWTAuth.ExpiresTime"),
			BufferTime:  v.GetInt64("JWTAuth.BufferTime"),
			Issuer:      v.GetString("JWTAuth.Issuer"),
		},
		UserModel:    query.Use(db),
		Rdb:          rdb,
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}
