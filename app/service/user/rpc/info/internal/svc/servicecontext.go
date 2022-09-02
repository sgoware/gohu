package svc

import (
	"github.com/hibiken/asynq"
	apollo "main/app/common/config"
	"main/app/common/log"
	CommentQuery "main/app/service/comment/dao/query"
	NotificationQuery "main/app/service/notification/dao/query"
	questionQuery "main/app/service/question/dao/query"
	UserQuery "main/app/service/user/dao/query"
	"main/app/service/user/rpc/info/internal/config"

	"github.com/go-redis/redis/v8"
)

type ServiceContext struct {
	Config config.Config

	UserModel         *UserQuery.Query
	QuestionModel     *questionQuery.Query
	CommentModel      *CommentQuery.Query
	NotificationModel *NotificationQuery.Query

	Rdb *redis.Client

	AsynqClient *asynq.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	logger := log.GetSugaredLogger()

	userDB, err := apollo.GetMysqlDB("user.yaml")
	if err != nil {
		logger.Fatalf("initialize mysql failed, err: %v", err)
	}

	questionDB, err := apollo.GetMysqlDB("question.yaml")
	if err != nil {
		logger.Fatalf("initialize mysql failed, err: %v", err)
	}

	commentDB, err := apollo.GetMysqlDB("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize mysql failed, err: %v", err)
	}

	notificationDB, err := apollo.GetMysqlDB("notification.yaml")
	if err != nil {
		logger.Fatalf("initialize mysql failed, err: %v", err)
	}

	rdb, err := apollo.GetRedisClient("user.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	return &ServiceContext{
		Config: c,

		UserModel:         UserQuery.Use(userDB),
		QuestionModel:     questionQuery.Use(questionDB),
		CommentModel:      CommentQuery.Use(commentDB),
		NotificationModel: NotificationQuery.Use(notificationDB),

		Rdb: rdb,

		AsynqClient: asynq.NewClient(c.AsynqClientConf),
	}
}
