package comment

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
	"main/app/service/comment/dao/model"
	"main/app/service/comment/dao/query"
	"main/app/service/comment/rpc/info/pb"
	"main/app/service/mq/asynq/processor/internal/config"
	"main/app/service/mq/asynq/processor/job"
	"main/app/utils/structx"
	"time"
)

type MsgCrudCommentSubjectHandler struct {
	Rdb          *redis.Client
	CommentModel *query.Query
	IdGenerator  *idgen.DefaultIdGenerator
}

type ScheduleUpdateCommentSubjectHandler struct {
	Rdb          *redis.Client
	CommentModel *query.Query
}

type ScheduleUpdateCommentIndexHandler struct {
	Rdb          *redis.Client
	CommentModel *query.Query
}

func NewMsgCrudCommentSubjectHandler(c config.Config) *MsgCrudCommentSubjectHandler {
	logger := log.GetSugaredLogger()

	rdb, err := apollo.GetRedisClient("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	commentDB, err := apollo.GetMysqlDB("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
	}

	idGenerator, err := apollo.NewIdGenerator("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize idGenerator failed, err: %v", err)
	}

	return &MsgCrudCommentSubjectHandler{
		Rdb:          rdb,
		CommentModel: query.Use(commentDB),
		IdGenerator:  idGenerator,
	}
}

func NewScheduleUpdateCommentSubjectHandler(c config.Config) *ScheduleUpdateCommentSubjectHandler {
	logger := log.GetSugaredLogger()

	rdb, err := apollo.GetRedisClient("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	commentDB, err := apollo.GetMysqlDB("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
	}

	return &ScheduleUpdateCommentSubjectHandler{
		Rdb:          rdb,
		CommentModel: query.Use(commentDB),
	}
}

func NewScheduleUpdateCommentIndexHandler(c config.Config) *ScheduleUpdateCommentIndexHandler {
	logger := log.GetSugaredLogger()

	rdb, err := apollo.GetRedisClient("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	commentDB, err := apollo.GetMysqlDB("comment.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
	}

	return &ScheduleUpdateCommentIndexHandler{
		Rdb:          rdb,
		CommentModel: query.Use(commentDB),
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

		commentSubjectProto := &pb.CommentSubject{
			Id:         commentSubjectId,
			ObjType:    payload.ObjType,
			ObjId:      payload.ObjId,
			Count:      payload.Count,
			RootCount:  payload.RootCount,
			State:      payload.State,
			Attrs:      payload.Attrs,
			CreateTime: payload.CreateTime.String(),
			UpdateTime: payload.UpdateTime.String(),
		}

		commentSubjectBytes, err := proto.Marshal(commentSubjectProto)
		if err != nil {
			return fmt.Errorf("marshal [commentSubjectBytes] failed, err: %v", err)
		}

		err = l.Rdb.Set(ctx,
			fmt.Sprintf("comment_subject_%d", payload.Id),
			commentSubjectBytes,
			time.Second*86400).Err()
		if err != nil {
			return fmt.Errorf("create [comment_subject] cache failed, err: %v", err)
		}

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

func (l *ScheduleUpdateCommentSubjectHandler) ProcessTask(ctx context.Context, _ *asynq.Task) (err error) {
	commentMembers, err := l.Rdb.SMembers(ctx,
		"comment_subject_comment_cnt_set").Result()
	if err != nil {
		return fmt.Errorf("get [comment_subject_comment_cnt_set] failed, err: %v", err)
	}
	err = l.Rdb.Del(ctx,
		"comment_subject_comment_cnt_set").Err()
	if err != nil {
		return fmt.Errorf("del [comment_subject_comment_cnt_set] failed, err: %v", err)
	}
	rootCommentMembers, err := l.Rdb.SMembers(ctx,
		"comment_subject_root_comment_cnt_set").Result()
	if err != nil {
		return fmt.Errorf("get [comment_subject_root_comment_cnt_set] failed, err: %v", err)
	}
	err = l.Rdb.Del(ctx,
		"comment_subject_root_comment_cnt_set").Err()
	if err != nil {
		return fmt.Errorf("del [comment_subject_root_comment_cnt_set] failed, err: %v", err)
	}

	commentSubjectModel := l.CommentModel.CommentSubject

	for _, commentMember := range commentMembers {
		commentSubjectId := cast.ToInt64(commentMember)
		commentCnt, err := l.Rdb.Get(ctx,
			fmt.Sprintf("comment_subject_comment_cnt_%d", commentSubjectId)).Int()
		if err != nil {
			return fmt.Errorf("get [comment_subject_comment_cnt] failed, err: %v", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("comment_subject_comment_cnt_%d", commentSubjectId)).Err()
		if err != nil {
			return fmt.Errorf("del [comment_subject_comment_cnt] failed, err: %v", err)
		}

		commentSubject, err := commentSubjectModel.WithContext(ctx).
			Select(commentSubjectModel.ID, commentSubjectModel.Count).
			Where(commentSubjectModel.ID.Eq(commentSubjectId)).
			First()
		if err != nil {
			return fmt.Errorf("query [comment_subject] failed, err: %v", err)
		}

		toCnt := commentSubject.Count + int32(commentCnt)

		_, err = commentSubjectModel.WithContext(ctx).
			Select(commentSubjectModel.ID, commentSubjectModel.Count).
			Where(commentSubjectModel.ID.Eq(commentSubjectId)).
			Update(commentSubjectModel.Count, toCnt)
		if err != nil {
			return fmt.Errorf("update [comment_subject] record failed, err: %v", err)
		}

		commentSubjectBytes, err := l.Rdb.Get(ctx,
			fmt.Sprintf("comment_subject_%d", commentSubjectId)).Bytes()
		if err == nil {
			commentSubjectProto := &pb.CommentSubject{}
			err = proto.Unmarshal(commentSubjectBytes, commentSubjectProto)
			if err != nil {
				return fmt.Errorf("unmarshal proto failed, err: %v", err)
			}
			commentSubjectProto.Count = toCnt
			commentSubjectBytes, err = proto.Marshal(commentSubjectProto)
			if err != nil {
				return fmt.Errorf("marshal [commentSubjectProto] failed, err: %v")
			}

			err = l.Rdb.Set(ctx,
				fmt.Sprintf("comment_subject_%d", commentSubjectId),
				commentSubjectBytes,
				time.Second*86400).Err()
			if err != nil {
				return fmt.Errorf("set [comment_subject] cache failed, err: %v", err)
			}
		} else {
			return fmt.Errorf("get [comment_subject] cache failed, err: %v", err)
		}
	}

	for _, rootCommentMember := range rootCommentMembers {
		commentSubjectId := cast.ToInt64(rootCommentMember)
		rootCommentCnt, err := l.Rdb.Get(ctx,
			fmt.Sprintf("comment_subject_root_comment_cnt_%d", commentSubjectId)).Int()
		if err != nil {
			return fmt.Errorf("get [comment_subject_root_comment_cnt] failed, err: %v", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("comment_subject_root_comment_cnt_%d", commentSubjectId)).Err()
		if err != nil {
			return fmt.Errorf("del [comment_subject_root_comment_cnt] failed, err: %v", err)
		}

		commentSubject, err := commentSubjectModel.WithContext(ctx).
			Select(commentSubjectModel.ID, commentSubjectModel.RootCount).
			Where(commentSubjectModel.ID.Eq(commentSubjectId)).
			First()
		if err != nil {
			return fmt.Errorf("query [comment_subject] failed, err: %v", err)
		}

		toCnt := commentSubject.RootCount + int32(rootCommentCnt)

		_, err = commentSubjectModel.WithContext(ctx).
			Select(commentSubjectModel.ID, commentSubjectModel.RootCount).
			Where(commentSubjectModel.ID.Eq(commentSubjectId)).
			Update(commentSubjectModel.RootCount, toCnt)
		if err != nil {
			return fmt.Errorf("update [comment_subject] record failed, err: %v", err)
		}

		commentSubjectBytes, err := l.Rdb.Get(ctx,
			fmt.Sprintf("comment_subject_%d", commentSubjectId)).Bytes()
		if err == nil {
			commentSubjectProto := &pb.CommentSubject{}
			err = proto.Unmarshal(commentSubjectBytes, commentSubjectProto)
			if err != nil {
				return fmt.Errorf("unmarshal proto failed, err: %v", err)
			}
			commentSubjectProto.RootCount = toCnt
			commentSubjectBytes, err = proto.Marshal(commentSubjectProto)
			if err != nil {
				return fmt.Errorf("marshal [commentSubjectProto] failed, err: %v")
			}

			err = l.Rdb.Set(ctx,
				fmt.Sprintf("comment_subject_%d", commentSubjectId),
				commentSubjectBytes,
				time.Second*86400).Err()
			if err != nil {
				return fmt.Errorf("set [comment_subject] cache failed, err: %v", err)
			}
		} else {
			return fmt.Errorf("get [comment_subject] cache failed, err: %v", err)
		}
	}
	return nil
}

func (l *ScheduleUpdateCommentIndexHandler) ProcessTask(ctx context.Context, _ *asynq.Task) (err error) {
	approveMembers, err := l.Rdb.SMembers(ctx,
		"comment_index_approve_cnt_set").Result()
	if err != nil {
		return fmt.Errorf("get [comment_index_approve_cnt_set] failed, err: %v", err)
	}
	err = l.Rdb.Del(ctx,
		"comment_index_approve_cnt_set").Err()
	if err != nil {
		return fmt.Errorf("del [comment_index_approve_cnt_set] failed ,err: %v", err)
	}

	commentIndexModel := l.CommentModel.CommentIndex

	for _, approveMember := range approveMembers {
		commentIndexId := cast.ToInt64(approveMember)

		approveCnt, err := l.Rdb.Get(ctx,
			fmt.Sprintf("comment_index_approve_cnt_%d", commentIndexId)).Int()
		if err != nil {
			return fmt.Errorf("get [comment_index_approve_cnt] failed, err: %d", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("comment_index_approve_cnt_%d", commentIndexId)).Err()
		if err != nil {
			return fmt.Errorf("del [comment_index_approve_cnt] failed, err: %v", err)
		}

		commentIndex, err := commentIndexModel.WithContext(ctx).
			Select(commentIndexModel.ID, commentIndexModel.ApproveCount).
			Where(commentIndexModel.ID.Eq(commentIndexId)).
			First()
		if err != nil {
			return fmt.Errorf("query [comment_index] record failed, err: %v", err)
		}

		toCnt := commentIndex.ApproveCount + int32(approveCnt)

		_, err = commentIndexModel.WithContext(ctx).
			Select(commentIndexModel.ID, commentIndexModel.ApproveCount).
			Where(commentIndexModel.ID.Eq(commentIndexId)).
			Update(commentIndexModel.ApproveCount, toCnt)
		if err != nil {
			return fmt.Errorf("update [comment_index] record failed, err: %v", err)
		}

		commentIndexBytes, err := l.Rdb.Get(ctx,
			fmt.Sprintf("comment_index_%d", commentIndexId)).Bytes()
		if err == nil {
			commentIndexProto := &pb.CommentIndex{}
			err = proto.Unmarshal(commentIndexBytes, commentIndexProto)
			if err != nil {
				return fmt.Errorf("unmarshal proto failed, err: %v", err)
			}
			commentIndexProto.ApproveCount = toCnt
			commentIndexBytes, err = proto.Marshal(commentIndexProto)
			if err != nil {
				return fmt.Errorf("marshal [commentIndexProto] failed, err: %v")
			}

			err = l.Rdb.Set(ctx,
				fmt.Sprintf("comment_index_%d", commentIndexId),
				commentIndexBytes,
				time.Second*86400).Err()
			if err != nil {
				return fmt.Errorf("set [comment_index] cache failed, err: %v", err)
			}
		}

	}

	return nil
}
