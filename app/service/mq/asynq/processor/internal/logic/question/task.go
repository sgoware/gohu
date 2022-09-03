package question

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/spf13/cast"
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

type ScheduleUpdateQuestionSubjectRecordHandler struct {
	Rdb           *redis.Client
	QuestionModel *query.Query
}

type ScheduleUpdateAnswerIndexRecordHandler struct {
	Rdb           *redis.Client
	QuestionModel *query.Query
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

func NewScheduleUpdateQuestionSubjectRecordHandler(c config.Config) *ScheduleUpdateQuestionSubjectRecordHandler {
	logger := log.GetSugaredLogger()

	rdb, err := apollo.GetRedisClient("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	questionDB, err := apollo.GetMysqlDB("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
	}

	return &ScheduleUpdateQuestionSubjectRecordHandler{
		Rdb:           rdb,
		QuestionModel: query.Use(questionDB),
	}
}

func NewScheduleUpdateAnswerIndexRecordHandler(c config.Config) *ScheduleUpdateAnswerIndexRecordHandler {
	logger := log.GetSugaredLogger()

	rdb, err := apollo.GetRedisClient("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	questionDB, err := apollo.GetMysqlDB("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
	}

	return &ScheduleUpdateAnswerIndexRecordHandler{
		Rdb:           rdb,
		QuestionModel: query.Use(questionDB),
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

func (l *ScheduleUpdateQuestionSubjectRecordHandler) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
	subCntMembers, err := l.Rdb.SMembers(ctx,
		"question_subject_sub_cnt_set").Result()
	if err != nil {
		return fmt.Errorf("get [question_subject_sub_cnt_set] failed ,err: %v", err)
	}
	err = l.Rdb.Del(ctx,
		fmt.Sprintf("question_subject_sub_cnt_set")).Err()
	if err != nil {
		return fmt.Errorf("delete [question_subject_sub_cnt_set] failed, err: %v", err)
	}

	answerCntMembers, err := l.Rdb.SMembers(ctx,
		"question_subject_answer_cnt_set").Result()
	if err != nil {
		return fmt.Errorf("get [question_subject_answer_cnt_set] failed ,err: %v", err)
	}
	err = l.Rdb.Del(ctx,
		fmt.Sprintf("question_subject_answer_cnt_set")).Err()
	if err != nil {
		return fmt.Errorf("delete [question_subject_answer_cnt_set] failed, err: %v", err)
	}

	viewCntMembers, err := l.Rdb.SMembers(ctx,
		"question_subject_view_cnt_set").Result()
	if err != nil {
		return fmt.Errorf("get [question_subject_view_cnt_set] failed ,err: %v", err)
	}
	err = l.Rdb.Del(ctx,
		fmt.Sprintf("question_subject_view_cnt_set")).Err()
	if err != nil {
		return fmt.Errorf("delete [question_subject_view_cnt_set] failed, err: %v", err)
	}

	questionSubjectModel := l.QuestionModel.QuestionSubject

	for _, subCntMember := range subCntMembers {
		subCnt, err := l.Rdb.Get(ctx,
			fmt.Sprintf("question_subejct_sub_cnt_%s", subCntMember)).Int()
		if err != nil {
			return fmt.Errorf("get [question_sub_cnt] failed, err: %d", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("question_subject_sub_cnt_%s", subCntMember)).Err()
		if err != nil {
			return fmt.Errorf("del [question_subject_sub_cnt] failed, err: %v", err)
		}

		questionSubject, err := questionSubjectModel.WithContext(ctx).
			Select(questionSubjectModel.ID, questionSubjectModel.SubCount).
			Where(questionSubjectModel.ID.Eq(cast.ToInt64(subCntMember))).
			First()
		if err != nil {
			return fmt.Errorf("query [question_subject] record failed, err: %v", err)
		}

		_, err = questionSubjectModel.WithContext(ctx).
			Select(questionSubjectModel.ID, questionSubjectModel.SubCount).
			Where(questionSubjectModel.ID.Eq(cast.ToInt64(subCntMember))).
			Update(questionSubjectModel.SubCount, int(questionSubject.SubCount)+subCnt)
		if err != nil {
			return fmt.Errorf("update [question_subejct] record failed, err: %v", err)
		}
	}

	for _, answerCntMember := range answerCntMembers {
		answerCnt, err := l.Rdb.Get(ctx,
			fmt.Sprintf("question_subject_answer_cnt_%s", answerCntMember)).Int()
		if err != nil {
			return fmt.Errorf("get [question_answer_cnt] failed, err: %d", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("question_subject_answer_cnt_%s", answerCntMember)).Err()
		if err != nil {
			return fmt.Errorf("del [question_subject_answer_cnt_] failed, err: %v", err)
		}

		questionSubject, err := questionSubjectModel.WithContext(ctx).
			Select(questionSubjectModel.ID, questionSubjectModel.AnswerCount).
			Where(questionSubjectModel.ID.Eq(cast.ToInt64(answerCntMember))).
			First()
		if err != nil {
			return fmt.Errorf("query [question_subject] record failed, err: %v", err)
		}

		_, err = questionSubjectModel.WithContext(ctx).
			Select(questionSubjectModel.ID, questionSubjectModel.AnswerCount).
			Where(questionSubjectModel.ID.Eq(cast.ToInt64(answerCntMember))).
			Update(questionSubjectModel.AnswerCount, int(questionSubject.AnswerCount)+answerCnt)
		if err != nil {
			return fmt.Errorf("update [question_subejct] record failed, err: %v", err)
		}
	}

	for _, viewCntMember := range viewCntMembers {
		viewCnt, err := l.Rdb.Get(ctx,
			fmt.Sprintf("question_subject_view_cnt_%s", viewCntMember)).Int()
		if err != nil {
			return fmt.Errorf("get [question_subject_view_cnt] failed, err: %d", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("question_subject_view_cnt_%s", viewCntMember)).Err()
		if err != nil {
			return fmt.Errorf("del [question_subject_view_cnt] failed, err: %v", err)
		}

		questionSubject, err := questionSubjectModel.WithContext(ctx).
			Select(questionSubjectModel.ID, questionSubjectModel.ViewCount).
			Where(questionSubjectModel.ID.Eq(cast.ToInt64(viewCntMember))).
			First()
		if err != nil {
			return fmt.Errorf("query [question_subject] record failed, err: %v", err)
		}

		_, err = questionSubjectModel.WithContext(ctx).
			Select(questionSubjectModel.ID, questionSubjectModel.ViewCount).
			Where(questionSubjectModel.ID.Eq(cast.ToInt64(viewCntMember))).
			Update(questionSubjectModel.ViewCount, int(questionSubject.ViewCount)+viewCnt)
		if err != nil {
			return fmt.Errorf("update [question_subejct] record failed, err: %v", err)
		}
	}

	return nil
}

func (l *ScheduleUpdateAnswerIndexRecordHandler) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
	approveMembers, err := l.Rdb.SMembers(ctx,
		"answer_index_approve_cnt_set").Result()
	if err != nil {
		return fmt.Errorf("get [answer_index_approve_cnt_set] failed ,err: %v", err)
	}
	err = l.Rdb.Del(ctx,
		fmt.Sprintf("answer_index_approve_cnt_set")).Err()
	if err != nil {
		return fmt.Errorf("delete [answer_index_approve_cnt_set] failed, err: %v", err)
	}

	likeMembers, err := l.Rdb.SMembers(ctx,
		"answer_index_like_cnt_set").Result()
	if err != nil {
		return fmt.Errorf("get [answer_index_like_cnt_set] failed ,err: %v", err)
	}
	err = l.Rdb.Del(ctx,
		fmt.Sprintf("answer_index_like_cnt_set")).Err()
	if err != nil {
		return fmt.Errorf("delete [answer_index_like_cnt_set] failed, err: %v", err)
	}

	collectMembers, err := l.Rdb.SMembers(ctx,
		"answer_index_collect_cnt_set").Result()
	if err != nil {
		return fmt.Errorf("get [answer_index_collect_cnt_set] failed ,err: %v", err)
	}
	err = l.Rdb.Del(ctx,
		fmt.Sprintf("answer_index_collect_cnt_set")).Err()
	if err != nil {
		return fmt.Errorf("delete [answer_index_collect_cnt_set] failed, err: %v", err)
	}

	answerIndexModel := l.QuestionModel.AnswerIndex

	for _, approveMember := range approveMembers {
		approveCnt, err := l.Rdb.Get(ctx,
			fmt.Sprintf("answer_index_approve_cnt_%s", approveMember)).Int()
		if err != nil {
			return fmt.Errorf("get [answer_index_approve_cnt] failed, err: %d", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("answer_index_approve_cnt_%s", approveMember)).Err()
		if err != nil {
			return fmt.Errorf("del [answer_index_approve_cnt] failed, err: %v", err)
		}

		answerIndex, err := answerIndexModel.WithContext(ctx).
			Select(answerIndexModel.ID, answerIndexModel.ApproveCount).
			Where(answerIndexModel.ID.Eq(cast.ToInt64(approveMember))).
			First()
		if err != nil {
			return fmt.Errorf("query [answer_index] record failed, err: %v", err)
		}

		_, err = answerIndexModel.WithContext(ctx).
			Select(answerIndexModel.ID, answerIndexModel.ApproveCount).
			Where(answerIndexModel.ID.Eq(cast.ToInt64(approveMember))).
			Update(answerIndexModel.ApproveCount, int(answerIndex.ApproveCount)+approveCnt)
		if err != nil {
			return fmt.Errorf("update [answer_index] record failed, err: %v", err)
		}
	}

	for _, likeMember := range likeMembers {
		likeCnt, err := l.Rdb.Get(ctx,
			fmt.Sprintf("answer_index_like_cnt_%s", likeMember)).Int()
		if err != nil {
			return fmt.Errorf("get [answer_index_like_cnt] failed, err: %d", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("answer_index_like_cnt_%s", likeMember)).Err()
		if err != nil {
			return fmt.Errorf("del [answer_index_like_cnt] failed, err: %v", err)
		}

		answerIndex, err := answerIndexModel.WithContext(ctx).
			Select(answerIndexModel.ID, answerIndexModel.LikeCount).
			Where(answerIndexModel.ID.Eq(cast.ToInt64(likeMember))).
			First()
		if err != nil {
			return fmt.Errorf("query [answer_index] record failed, err: %v", err)
		}

		_, err = answerIndexModel.WithContext(ctx).
			Select(answerIndexModel.ID, answerIndexModel.LikeCount).
			Where(answerIndexModel.ID.Eq(cast.ToInt64(likeMember))).
			Update(answerIndexModel.LikeCount, int(answerIndex.LikeCount)+likeCnt)
		if err != nil {
			return fmt.Errorf("update [answer_index] record failed, err: %v", err)
		}
	}

	for _, collectMember := range collectMembers {
		collectCnt, err := l.Rdb.Get(ctx,
			fmt.Sprintf("answer_index_collect_cnt_%s", collectMember)).Int()
		if err != nil {
			return fmt.Errorf("get [answer_index_collect_cnt] failed, err: %d", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("answer_index_collect_cnt_%s", collectMember)).Err()
		if err != nil {
			return fmt.Errorf("del [answer_index_collect_cnt] failed, err: %v", err)
		}

		answerIndex, err := answerIndexModel.WithContext(ctx).
			Select(answerIndexModel.ID, answerIndexModel.CollectCount).
			Where(answerIndexModel.ID.Eq(cast.ToInt64(collectMember))).
			First()
		if err != nil {
			return fmt.Errorf("query [answer_index] record failed, err: %v", err)
		}

		_, err = answerIndexModel.WithContext(ctx).
			Select(answerIndexModel.ID, answerIndexModel.CollectCount).
			Where(answerIndexModel.ID.Eq(cast.ToInt64(collectMember))).
			Update(answerIndexModel.CollectCount, int(answerIndex.CollectCount)+collectCnt)
		if err != nil {
			return fmt.Errorf("update [answer_index] record failed, err: %v", err)
		}
	}

	return nil
}
