package svc

import (
	"github.com/go-redis/redis/v8"
	apollo "main/app/common/config"
	"main/app/common/log"
	commentQuery "main/app/service/comment/dao/query"
	commentRpc "main/app/service/comment/rpc/crud/crud"
	"main/app/service/mq/asynq/processor/internal/config"
	questionQuery "main/app/service/question/dao/query"
	questionRpc "main/app/service/question/rpc/crud/crud"
	userQuery "main/app/service/user/dao/query"
	userRpc "main/app/service/user/rpc/crud/crud"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	AsynqServer *asynq.Server

	UserCrudRpcClient     userRpc.Crud
	QuestionCrudRpcClient questionRpc.Crud
	CommentCrudRpcClient  commentRpc.Crud

	Rdb *redis.Client

	UserModel     *userQuery.Query
	QuestionModel *questionQuery.Query
	CommentModel  *commentQuery.Query
}

func NewServiceContext(c config.Config) *ServiceContext {
	logger := log.GetSugaredLogger()

	rdb, err := apollo.GetRedisClient("user.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	userDB, err := apollo.GetMysqlDB("user.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
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

		UserCrudRpcClient:     userRpc.NewCrud(zrpc.MustNewClient(c.UserCrudRpcClientConf)),
		QuestionCrudRpcClient: questionRpc.NewCrud(zrpc.MustNewClient(c.QuestionCrudRpcClientConf)),
		CommentCrudRpcClient:  commentRpc.NewCrud(zrpc.MustNewClient(c.CommentCrudRpcClientConf)),

		Rdb: rdb,

		UserModel:     userQuery.Use(userDB),
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
