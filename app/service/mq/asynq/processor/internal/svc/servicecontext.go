package svc

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/zrpc"
	"main/app/common/log"
	comment "main/app/service/comment/rpc/crud/crud"
	"main/app/service/mq/asynq/processor/internal/config"
	question "main/app/service/question/rpc/crud/crud"
	user "main/app/service/user/rpc/crud/crud"
)

type ServiceContext struct {
	Config config.Config

	AsynqServer *asynq.Server

	UserCrudRpcClient     user.Crud
	QuestionCrudRpcClient question.Crud
	CommentCrudRpcClient  comment.Crud
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: config.Config{},

		AsynqServer: newAsynqServer(c),

		UserCrudRpcClient:     user.NewCrud(zrpc.MustNewClient(c.UserCrudRpcClientConf)),
		QuestionCrudRpcClient: question.NewCrud(zrpc.MustNewClient(c.QuestionCrudRpcClientConf)),
		CommentCrudRpcClient:  comment.NewCrud(zrpc.MustNewClient(c.CommentCrudRpcClientConf)),
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
