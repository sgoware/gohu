package question

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/spf13/cast"
	"github.com/yitter/idgenerator-go/idgen"
	"google.golang.org/protobuf/proto"
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/service/mq/asynq/processor/internal/config"
	"main/app/service/mq/asynq/processor/job"
	"main/app/service/question/dao/model"
	"main/app/service/question/dao/pb"
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

type ScheduleUpdateQuestionSubjectHandler struct {
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

func NewScheduleUpdateQuestionSubjectHandler(c config.Config) *ScheduleUpdateQuestionSubjectHandler {
	logger := log.GetSugaredLogger()

	rdb, err := apollo.GetRedisClient("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	questionDB, err := apollo.GetMysqlDB("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
	}

	return &ScheduleUpdateQuestionSubjectHandler{
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

func (l *ScheduleUpdateQuestionSubjectHandler) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
	subCntMembers, err := l.Rdb.SMembers(ctx,
		"question_subject_sub_cnt_set").Result()
	if err != nil {
		return fmt.Errorf("get [question_subject_sub_cnt_set] failed ,err: %v", err)
	}
	err = l.Rdb.Del(ctx,
		"question_subject_sub_cnt_set").Err()
	if err != nil {
		return fmt.Errorf("delete [question_subject_sub_cnt_set] failed, err: %v", err)
	}

	answerCntMembers, err := l.Rdb.SMembers(ctx,
		"question_subject_answer_cnt_set").Result()
	if err != nil {
		return fmt.Errorf("get [question_subject_answer_cnt_set] failed ,err: %v", err)
	}
	err = l.Rdb.Del(ctx,
		"question_subject_answer_cnt_set").Err()
	if err != nil {
		return fmt.Errorf("delete [question_subject_answer_cnt_set] failed, err: %v", err)
	}

	viewCntMembers, err := l.Rdb.SMembers(ctx,
		"question_subject_view_cnt_set").Result()
	if err != nil {
		return fmt.Errorf("get [question_subject_view_cnt_set] failed ,err: %v", err)
	}
	err = l.Rdb.Del(ctx,
		"question_subject_view_cnt_set").Err()
	if err != nil {
		return fmt.Errorf("delete [question_subject_view_cnt_set] failed, err: %v", err)
	}

	questionSubjectModel := l.QuestionModel.QuestionSubject

	for _, subCntMember := range subCntMembers {
		questionSubjectId := cast.ToInt64(subCntMember)
		subCnt, err := l.Rdb.Get(ctx,
			fmt.Sprintf("question_subejct_sub_cnt_%d", questionSubjectId)).Int()
		if err != nil {
			return fmt.Errorf("get [question_sub_cnt] failed, err: %d", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("question_subject_sub_cnt_%d", questionSubjectId)).Err()
		if err != nil {
			return fmt.Errorf("del [question_subject_sub_cnt] failed, err: %v", err)
		}

		questionSubject, err := questionSubjectModel.WithContext(ctx).
			Select(questionSubjectModel.ID, questionSubjectModel.SubCount).
			Where(questionSubjectModel.ID.Eq(questionSubjectId)).
			First()
		if err != nil {
			return fmt.Errorf("query [question_subject] record failed, err: %v", err)
		}

		toCnt := questionSubject.SubCount + int32(subCnt)

		_, err = questionSubjectModel.WithContext(ctx).
			Select(questionSubjectModel.ID, questionSubjectModel.SubCount).
			Where(questionSubjectModel.ID.Eq(questionSubjectId)).
			Update(questionSubjectModel.SubCount, toCnt)
		if err != nil {
			return fmt.Errorf("update [question_subejct] record failed, err: %v", err)
		}

		questionSubjectBytes, err := l.Rdb.Get(ctx,
			fmt.Sprintf("question_subject_%d", questionSubjectId)).Bytes()
		if err == nil {
			questionSubjectProto := &pb.QuestionSubject{}
			err = proto.Unmarshal(questionSubjectBytes, questionSubjectProto)
			if err != nil {
				return fmt.Errorf("unmarshal proto failed, err: %v", err)
			}
			questionSubjectProto.SubCount = toCnt
			questionSubjectBytes, err = proto.Marshal(questionSubjectProto)
			if err != nil {
				return fmt.Errorf("marshal [questionSubjectProto] failed, err: %v")
			}
		} else {
			return fmt.Errorf("get [question_subject] cache failed, err: %v", err)
		}
	}

	for _, answerCntMember := range answerCntMembers {
		questionSubjectId := cast.ToInt64(answerCntMember)
		answerCnt, err := l.Rdb.Get(ctx,
			fmt.Sprintf("question_subject_answer_cnt_%d", questionSubjectId)).Int()
		if err != nil {
			return fmt.Errorf("get [question_answer_cnt] failed, err: %d", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("question_subject_answer_cnt_%d", questionSubjectId)).Err()
		if err != nil {
			return fmt.Errorf("del [question_subject_answer_cnt_] failed, err: %v", err)
		}

		questionSubject, err := questionSubjectModel.WithContext(ctx).
			Select(questionSubjectModel.ID, questionSubjectModel.AnswerCount).
			Where(questionSubjectModel.ID.Eq(questionSubjectId)).
			First()
		if err != nil {
			return fmt.Errorf("query [question_subject] record failed, err: %v", err)
		}

		toCnt := questionSubject.AnswerCount + int32(answerCnt)

		_, err = questionSubjectModel.WithContext(ctx).
			Select(questionSubjectModel.ID, questionSubjectModel.AnswerCount).
			Where(questionSubjectModel.ID.Eq(questionSubjectId)).
			Update(questionSubjectModel.AnswerCount, toCnt)
		if err != nil {
			return fmt.Errorf("update [question_subejct] record failed, err: %v", err)
		}

		questionSubjectBytes, err := l.Rdb.Get(ctx,
			fmt.Sprintf("question_subject_%d", questionSubjectId)).Bytes()
		if err == nil {
			questionSubjectProto := &pb.QuestionSubject{}
			err = proto.Unmarshal(questionSubjectBytes, questionSubjectProto)
			if err != nil {
				return fmt.Errorf("unmarshal proto failed, err: %v", err)
			}
			questionSubjectProto.AnswerCount = toCnt
			questionSubjectBytes, err = proto.Marshal(questionSubjectProto)
			if err != nil {
				return fmt.Errorf("marshal [questionSubjectProto] failed, err: %v")
			}
		} else {
			return fmt.Errorf("get [question_subject] cache failed, err: %v", err)
		}
	}

	for _, viewCntMember := range viewCntMembers {
		questionSubjectId := cast.ToInt64(viewCntMember)
		viewCnt, err := l.Rdb.Get(ctx,
			fmt.Sprintf("question_subject_view_cnt_%d", questionSubjectId)).Int()
		if err != nil {
			return fmt.Errorf("get [question_subject_view_cnt] failed, err: %d", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("question_subject_view_cnt_%d", questionSubjectId)).Err()
		if err != nil {
			return fmt.Errorf("del [question_subject_view_cnt] failed, err: %v", err)
		}

		questionSubject, err := questionSubjectModel.WithContext(ctx).
			Select(questionSubjectModel.ID, questionSubjectModel.ViewCount).
			Where(questionSubjectModel.ID.Eq(questionSubjectId)).
			First()
		if err != nil {
			return fmt.Errorf("query [question_subject] record failed, err: %v", err)
		}

		toCnt := questionSubject.ViewCount + int64(viewCnt)

		_, err = questionSubjectModel.WithContext(ctx).
			Select(questionSubjectModel.ID, questionSubjectModel.ViewCount).
			Where(questionSubjectModel.ID.Eq(questionSubjectId)).
			Update(questionSubjectModel.ViewCount, toCnt)
		if err != nil {
			return fmt.Errorf("update [question_subejct] record failed, err: %v", err)
		}

		questionSubjectBytes, err := l.Rdb.Get(ctx,
			fmt.Sprintf("question_subject_%d", questionSubjectId)).Bytes()
		if err == nil {
			questionSubjectProto := &pb.QuestionSubject{}
			err = proto.Unmarshal(questionSubjectBytes, questionSubjectProto)
			if err != nil {
				return fmt.Errorf("unmarshal proto failed, err: %v", err)
			}
			questionSubjectProto.ViewCount = toCnt
			questionSubjectBytes, err = proto.Marshal(questionSubjectProto)
			if err != nil {
				return fmt.Errorf("marshal [questionSubjectProto] failed, err: %v")
			}
		} else {
			return fmt.Errorf("get [question_subject] cache failed, err: %v", err)
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
		"answer_index_approve_cnt_set").Err()
	if err != nil {
		return fmt.Errorf("delete [answer_index_approve_cnt_set] failed, err: %v", err)
	}

	likeMembers, err := l.Rdb.SMembers(ctx,
		"answer_index_like_cnt_set").Result()
	if err != nil {
		return fmt.Errorf("get [answer_index_like_cnt_set] failed ,err: %v", err)
	}
	err = l.Rdb.Del(ctx,
		"answer_index_like_cnt_set").Err()
	if err != nil {
		return fmt.Errorf("delete [answer_index_like_cnt_set] failed, err: %v", err)
	}

	collectMembers, err := l.Rdb.SMembers(ctx,
		"answer_index_collect_cnt_set").Result()
	if err != nil {
		return fmt.Errorf("get [answer_index_collect_cnt_set] failed ,err: %v", err)
	}
	err = l.Rdb.Del(ctx,
		"answer_index_collect_cnt_set").Err()
	if err != nil {
		return fmt.Errorf("delete [answer_index_collect_cnt_set] failed, err: %v", err)
	}

	answerIndexModel := l.QuestionModel.AnswerIndex

	for _, approveMember := range approveMembers {
		answerIndexId := cast.ToInt64(approveMember)
		approveCnt, err := l.Rdb.Get(ctx,
			fmt.Sprintf("answer_index_approve_cnt_%d", answerIndexId)).Int()
		if err != nil {
			return fmt.Errorf("get [answer_index_approve_cnt] failed, err: %d", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("answer_index_approve_cnt_%d", answerIndexId)).Err()
		if err != nil {
			return fmt.Errorf("del [answer_index_approve_cnt] failed, err: %v", err)
		}

		answerIndex, err := answerIndexModel.WithContext(ctx).
			Select(answerIndexModel.ID, answerIndexModel.ApproveCount).
			Where(answerIndexModel.ID.Eq(answerIndexId)).
			First()
		if err != nil {
			return fmt.Errorf("query [answer_index] record failed, err: %v", err)
		}

		toCnt := answerIndex.ApproveCount + int32(approveCnt)

		_, err = answerIndexModel.WithContext(ctx).
			Select(answerIndexModel.ID, answerIndexModel.ApproveCount).
			Where(answerIndexModel.ID.Eq(answerIndexId)).
			Update(answerIndexModel.ApproveCount, toCnt)
		if err != nil {
			return fmt.Errorf("update [answer_index] record failed, err: %v", err)
		}

		answerIndexBytes, err := l.Rdb.Get(ctx,
			fmt.Sprintf("answer_index_%d", answerIndexId)).Bytes()
		if err == nil {
			answerIndexProto := &pb.AnswerIndex{}
			err = proto.Unmarshal(answerIndexBytes, answerIndexProto)
			if err != nil {
				return fmt.Errorf("unmarshal proto failed, err: %v", err)
			}
			answerIndexProto.ApproveCount = toCnt
			answerIndexBytes, err = proto.Marshal(answerIndexProto)
			if err != nil {
				return fmt.Errorf("marshal [answerIndexProto] failed, err: %v")
			}
		} else {
			return fmt.Errorf("get [answer_index] cache failed, err: %v", err)
		}
	}

	for _, likeMember := range likeMembers {
		answerIndexId := cast.ToInt64(likeMember)
		likeCnt, err := l.Rdb.Get(ctx,
			fmt.Sprintf("answer_index_like_cnt_%d", answerIndexId)).Int()
		if err != nil {
			return fmt.Errorf("get [answer_index_like_cnt] failed, err: %d", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("answer_index_like_cnt_%d", answerIndexId)).Err()
		if err != nil {
			return fmt.Errorf("del [answer_index_like_cnt] failed, err: %v", err)
		}

		answerIndex, err := answerIndexModel.WithContext(ctx).
			Select(answerIndexModel.ID, answerIndexModel.LikeCount).
			Where(answerIndexModel.ID.Eq(answerIndexId)).
			First()
		if err != nil {
			return fmt.Errorf("query [answer_index] record failed, err: %v", err)
		}

		toCnt := answerIndex.LikeCount + int32(likeCnt)

		_, err = answerIndexModel.WithContext(ctx).
			Select(answerIndexModel.ID, answerIndexModel.LikeCount).
			Where(answerIndexModel.ID.Eq(answerIndexId)).
			Update(answerIndexModel.LikeCount, toCnt)
		if err != nil {
			return fmt.Errorf("update [answer_index] record failed, err: %v", err)
		}

		answerIndexBytes, err := l.Rdb.Get(ctx,
			fmt.Sprintf("answer_index_%d", answerIndexId)).Bytes()
		if err == nil {
			answerIndexProto := &pb.AnswerIndex{}
			err = proto.Unmarshal(answerIndexBytes, answerIndexProto)
			if err != nil {
				return fmt.Errorf("unmarshal proto failed, err: %v", err)
			}
			answerIndexProto.LikeCount = toCnt
			answerIndexBytes, err = proto.Marshal(answerIndexProto)
			if err != nil {
				return fmt.Errorf("marshal [answerIndexProto] failed, err: %v")
			}
		} else {
			return fmt.Errorf("get [answer_index] cache failed, err: %v", err)
		}
	}

	for _, collectMember := range collectMembers {
		answerIndexId := cast.ToInt64(collectMember)
		collectCnt, err := l.Rdb.Get(ctx,
			fmt.Sprintf("answer_index_collect_cnt_%d", answerIndexId)).Int()
		if err != nil {
			return fmt.Errorf("get [answer_index_collect_cnt] failed, err: %d", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("answer_index_collect_cnt_%d", answerIndexId)).Err()
		if err != nil {
			return fmt.Errorf("del [answer_index_collect_cnt] failed, err: %v", err)
		}

		answerIndex, err := answerIndexModel.WithContext(ctx).
			Select(answerIndexModel.ID, answerIndexModel.CollectCount).
			Where(answerIndexModel.ID.Eq(answerIndexId)).
			First()
		if err != nil {
			return fmt.Errorf("query [answer_index] record failed, err: %v", err)
		}

		toCnt := answerIndex.CollectCount + int32(collectCnt)

		_, err = answerIndexModel.WithContext(ctx).
			Select(answerIndexModel.ID, answerIndexModel.CollectCount).
			Where(answerIndexModel.ID.Eq(answerIndexId)).
			Update(answerIndexModel.CollectCount, toCnt)
		if err != nil {
			return fmt.Errorf("update [answer_index] record failed, err: %v", err)
		}

		answerIndexBytes, err := l.Rdb.Get(ctx,
			fmt.Sprintf("answer_index_%d", answerIndexId)).Bytes()
		if err == nil {
			answerIndexProto := &pb.AnswerIndex{}
			err = proto.Unmarshal(answerIndexBytes, answerIndexProto)
			if err != nil {
				return fmt.Errorf("unmarshal proto failed, err: %v", err)
			}
			answerIndexProto.CollectCount = toCnt
			answerIndexBytes, err = proto.Marshal(answerIndexProto)
			if err != nil {
				return fmt.Errorf("marshal [answerIndexProto] failed, err: %v")
			}
		} else {
			return fmt.Errorf("get [answer_index] cache failed, err: %v", err)
		}
	}

	return nil
}
