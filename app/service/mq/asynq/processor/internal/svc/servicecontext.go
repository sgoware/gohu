package svc

import (
	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	apollo "main/app/common/config"
	"main/app/common/log"
	commentQuery "main/app/service/comment/dao/query"
	"main/app/service/mq/asynq/processor/internal/config"
	questionQuery "main/app/service/question/dao/query"
)

type ServiceContext struct {
	Config config.Config

	AsynqServer *asynq.Server

	Rdb *redis.Client

	QuestionModel *questionQuery.Query
	CommentModel  *commentQuery.Query
}

func NewServiceContext(c config.Config) *ServiceContext {
	logger := log.GetSugaredLogger()

	rdb, err := apollo.GetRedisClient("user.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	questionDB, err := apollo.GetMysqlDB("question.yaml")
	if err != nil {
		logger.Fatalf("initialize question mysql failed, err: %v", err)
	}

	commentDB, err := apollo.GetMysqlDB("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize comment mysql failed ,err: %v", err)
	}

	return &ServiceContext{
		Config: config.Config{},

		AsynqServer: newAsynqServer(c),

		Rdb: rdb,

		QuestionModel: questionQuery.Use(questionDB),
		CommentModel:  commentQuery.Use(commentDB),
	}
}

func newAsynqServer(c config.Config) *asynq.Server {
	logger := log.GetSugaredLogger()
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: c.Redis.Host, Password: c.Redis.Pass},
		asynq.Config{
			Logger:      logger,
			LogLevel:    asynq.DebugLevel,
			Concurrency: 20, //max concurrent process job task num
		},
	)
}
