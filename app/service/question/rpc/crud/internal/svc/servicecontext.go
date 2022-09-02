package svc

import (
	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/yitter/idgenerator-go/idgen"
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/service/question/dao/query"
	"main/app/service/question/rpc/crud/internal/config"
)

type ServiceContext struct {
	Config config.Config

	QuestionModel *query.Query
	Rdb           *redis.Client

	AsynqClient *asynq.Client

	IdGenerator *idgen.DefaultIdGenerator
}

func NewServiceContext(c config.Config) *ServiceContext {
	logger := log.GetSugaredLogger()

	db, err := apollo.GetMysqlDB("question.yaml")
	if err != nil {
		logger.Fatalf("initialize mysql failed, err: %v", err)
	}

	rdb, err := apollo.GetRedisClient("question.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	idGenerator, err := apollo.NewIdGenerator("question.yaml")
	if err != nil {
		logger.Fatalf("initialize idGenerator failed, err: %v", err)
	}

	return &ServiceContext{
		Config: c,

		QuestionModel: query.Use(db),
		Rdb:           rdb,

		AsynqClient: asynq.NewClient(c.AsynqClientConf),

		IdGenerator: idGenerator,
	}
}
