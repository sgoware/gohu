package question

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/yitter/idgenerator-go/idgen"
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/service/mq/asynq/processor/internal/config"
	"main/app/service/mq/asynq/processor/job"
	"main/app/service/question/dao/model"
	"main/app/service/question/dao/query"
	"main/app/utils/structx"
)

type MsgCrudQuestionSubjectHandler struct {
	Rdb           *redis.Client
	QuestionModel *query.Query
	IdGenerator   *idgen.DefaultIdGenerator
}

type MsgCrudQuestionContentHandler struct {
	Rdb           *redis.Client
	QuestionModel *query.Query
	IdGenerator   *idgen.DefaultIdGenerator
}

func NewMsgCrudQuestionSubjectHandler(c config.Config) *MsgCrudQuestionSubjectHandler {
	logger := log.GetSugaredLogger()

	rdb, err := apollo.GetRedisClient("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	questionDB, err := apollo.GetMysqlDB("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
	}

	idGenerator, err := apollo.NewIdGenerator("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize idGenerator failed, err: %v", err)
	}

	return &MsgCrudQuestionSubjectHandler{
		Rdb:           rdb,
		QuestionModel: query.Use(questionDB),
		IdGenerator:   idGenerator,
	}
}

func NewMsgCrudQuestionContentHandler(c config.Config) *MsgCrudQuestionContentHandler {
	logger := log.GetSugaredLogger()

	rdb, err := apollo.GetRedisClient("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	questionDB, err := apollo.GetMysqlDB("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
	}

	idGenerator, err := apollo.NewIdGenerator("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize idGenerator failed, err: %v", err)
	}

	return &MsgCrudQuestionContentHandler{
		Rdb:           rdb,
		QuestionModel: query.Use(questionDB),
		IdGenerator:   idGenerator,
	}
}

func (l *MsgCrudQuestionSubjectHandler) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
	var payload job.MsgCrudQuestionSubjectRecordPayload
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal [MsgCrudQuestionSubjectRecordPayload] failed ,err: %v", err)
	}

	questionSubjectModel := l.QuestionModel.QuestionSubject

	switch payload.Action {
	case 1:
		err = questionSubjectModel.WithContext(ctx).
			Create(&model.QuestionSubject{
				ID:          payload.Id,
				UserID:      payload.UserId,
				IPLoc:       payload.IpLoc,
				Title:       payload.Title,
				Topic:       payload.Topic,
				Tag:         payload.Tag,
				SubCount:    payload.SubCount,
				AnswerCount: payload.AnswerCount,
				ViewCount:   payload.ViewCount,
				State:       payload.State,
				Attrs:       payload.Attrs,
				CreateTime:  payload.CreateTime,
				UpdateTime:  payload.UpdateTime,
			})
		if err != nil {
			return fmt.Errorf("create [question_subject] record failed, err: %v", err)
		}

	case 2:
		questionSubject := &model.QuestionSubject{}

		err = structx.SyncWithNoZero(payload, questionSubject)
		if err != nil {
			return fmt.Errorf("sync struct [questionSubject] failed, err: %v", err)
		}

		_, err = questionSubjectModel.WithContext(ctx).
			Where(questionSubjectModel.ID.Eq(questionSubject.ID)).
			Updates(questionSubject)
		if err != nil {
			return fmt.Errorf("update [question_subject] record failed, err: %v", err)
		}

	case 3:
		_, err = questionSubjectModel.WithContext(ctx).
			Where(questionSubjectModel.ID.Eq(payload.Id)).
			Delete()
		if err != nil {
			return fmt.Errorf("delete [question_subject] record failed, err: %v", err)
		}
	}

	return nil
}

func (l *MsgCrudQuestionContentHandler) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
	var payload job.MsgCrudQuestionContentRecordPayload
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal [MsgCrudQuestionSubjectRecordPayload] failed ,err: %v", err)
	}

	questionContentModel := l.QuestionModel.QuestionContent

	switch payload.Action {
	case 1:
		err = questionContentModel.WithContext(ctx).
			Create(&model.QuestionContent{
				QuestionID: payload.QuestionId,
				Content:    payload.Content,
				Meta:       payload.Meta,
				CreateTime: payload.CreateTime,
				UpdateTime: payload.UpdateTime,
			})
		if err != nil {
			return fmt.Errorf("create [question_Content] record failed, err: %v", err)
		}

	case 2:
		questionContent := &model.QuestionContent{}

		err = structx.SyncWithNoZero(payload, questionContent)
		if err != nil {
			return fmt.Errorf("sync struct [questionContent] failed, err: %v", err)
		}

		_, err = questionContentModel.WithContext(ctx).
			Where(questionContentModel.QuestionID.Eq(payload.QuestionId)).
			Updates(questionContent)
		if err != nil {
			return fmt.Errorf("update [question_Content] record failed, err: %v", err)
		}

	case 3:
		_, err = questionContentModel.WithContext(ctx).
			Where(questionContentModel.QuestionID.Eq(payload.QuestionId)).
			Delete()
		if err != nil {
			return fmt.Errorf("delete [question_Content] record failed, err: %v", err)
		}
	}

	return nil
}
