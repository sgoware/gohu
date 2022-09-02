package comment

import (
	"github.com/go-redis/redis/v8"
	"github.com/yitter/idgenerator-go/idgen"
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/service/comment/dao/query"
	"main/app/service/question/mq/config"
)

type MsgCrudCommentSubjectHandler struct {
	Rdb          *redis.Client
	CommentModel *query.Query
	IdGenerator  *idgen.DefaultIdGenerator
}

func NewMsgCrudCommentSubjectHandler(c config.Config) *MsgCrudCommentSubjectHandler {
	logger := log.GetSugaredLogger()

	rdb, err := apollo.GetRedisClient("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	userDB, err := apollo.GetMysqlDB("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
	}

	idGenerator, err := apollo.NewIdGenerator("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize idGenerator failed, err: %v", err)
	}

	return &MsgCrudCommentSubjectHandler{
		Rdb:          rdb,
		CommentModel: query.Use(userDB),
		IdGenerator:  idGenerator,
	}
}
