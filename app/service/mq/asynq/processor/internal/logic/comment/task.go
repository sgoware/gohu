package comment

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/yitter/idgenerator-go/idgen"
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/service/comment/dao/model"
	"main/app/service/comment/dao/query"
	"main/app/service/mq/asynq/processor/job"
	"main/app/service/question/mq/config"
	"main/app/utils/structx"
	"time"
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

func (l *MsgCrudCommentSubjectHandler) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
	var payload job.MsgCrudCommentSubjectPayload
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal MsgCrudCommentSubjectPayload failed, err: %v", err)
	}

	commentSubjectModel := l.CommentModel.CommentSubject

	nowTime := time.Now()

	switch payload.Action {
	case 1:
		// 创建
		commentSubjectId := l.IdGenerator.NewLong()

		err = commentSubjectModel.WithContext(ctx).
			Create(&model.CommentSubject{
				ID:         commentSubjectId,
				ObjType:    payload.ObjType,
				ObjID:      payload.ObjId,
				Count:      payload.Count,
				RootCount:  payload.RootCount,
				State:      payload.State,
				Attrs:      payload.Attrs,
				CreateTime: nowTime,
				UpdateTime: nowTime,
			})
		if err != nil {
			return fmt.Errorf("create [comment_subject] record failed, err: %v", err)
		}

	case 2:
		// 更新
		commentSubject := &model.CommentSubject{}

		err = structx.SyncWithNoZero(payload, commentSubject)
		if err != nil {
			return fmt.Errorf("sync struct [CommentSubject] failed, err: %v", err)
		}

		_, err = commentSubjectModel.WithContext(ctx).
			Where(commentSubjectModel.ID.Eq(payload.Id)).
			Updates(commentSubject)
		if err != nil {
			return fmt.Errorf("update [comment_subject] record failed, err: %v", err)
		}

	case 3:
		// 删除
		_, err = commentSubjectModel.WithContext(ctx).
			Where(commentSubjectModel.ID.Eq(payload.Id)).
			Delete()
		if err != nil {
			return fmt.Errorf("delete [comment_subject] record failed, err: %v", err)
		}
	}
	return nil
}
