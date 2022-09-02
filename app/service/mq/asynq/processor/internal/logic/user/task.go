package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/proto"
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/service/mq/asynq/processor/internal/config"
	"main/app/service/mq/asynq/processor/job"
	"main/app/service/user/dao/model"
	"main/app/service/user/dao/pb"
	"main/app/service/user/dao/query"
	"main/app/utils/structx"
	"time"
)

type MsgCreateUserSubjectHandler struct {
	Rdb       *redis.Client
	UserModel *query.Query
}

type MsgUpdateUserSubjectRecordHandler struct {
	Rdb       *redis.Client
	UserModel *query.Query
}

type MsgUpdateUserSubjectCacheHandler struct {
	Rdb       *redis.Client
	UserModel *query.Query
}

type MsgAddUserSubjectCacheHandler struct {
	Rdb       *redis.Client
	UserModel *query.Query
}

type ScheduleUpdateUserSubjectRecordHandler struct {
	Rdb       *redis.Client
	UserModel *query.Query
}

type MsgUpdateUserCollectCacheHandler struct {
	Rdb       *redis.Client
	UserModel *query.Query
}

type ScheduleUpdateUserCollectRecordHandler struct {
	Rdb       *redis.Client
	UserModel *query.Query
}

func NewCreateUserSubjectRecordHandler(c config.Config) *MsgCreateUserSubjectHandler {
	logger := log.GetSugaredLogger()

	rdb, err := apollo.GetRedisClient("user.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	userDB, err := apollo.GetMysqlDB("user.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
	}

	return &MsgCreateUserSubjectHandler{
		Rdb: rdb,

		UserModel: query.Use(userDB),
	}
}

func NewUpdateUserSubjectRecordHandler(c config.Config) *MsgUpdateUserSubjectRecordHandler {
	logger := log.GetSugaredLogger()

	rdb, err := apollo.GetRedisClient("user.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	userDB, err := apollo.GetMysqlDB("user.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
	}

	return &MsgUpdateUserSubjectRecordHandler{
		Rdb: rdb,

		UserModel: query.Use(userDB),
	}
}

func NewUpdateUserSubjectCacheHandler(c config.Config) *MsgUpdateUserSubjectCacheHandler {
	logger := log.GetSugaredLogger()
	rdb, err := apollo.GetRedisClient("user.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	userDB, err := apollo.GetMysqlDB("user.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
	}

	return &MsgUpdateUserSubjectCacheHandler{
		Rdb: rdb,

		UserModel: query.Use(userDB),
	}
}

func NewMsgAddUserSubjectCacheHandler(c config.Config) *MsgAddUserSubjectCacheHandler {
	logger := log.GetSugaredLogger()

	rdb, err := apollo.GetRedisClient("user.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	userDB, err := apollo.GetMysqlDB("user.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
	}

	return &MsgAddUserSubjectCacheHandler{
		Rdb: rdb,

		UserModel: query.Use(userDB),
	}
}

func NewScheduleUpdateUserSubjectRecordHandler(c config.Config) *ScheduleUpdateUserSubjectRecordHandler {
	logger := log.GetSugaredLogger()

	rdb, err := apollo.GetRedisClient("user.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	userDB, err := apollo.GetMysqlDB("user.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
	}

	return &ScheduleUpdateUserSubjectRecordHandler{
		Rdb: rdb,

		UserModel: query.Use(userDB),
	}
}

func NewMsgUpdateUserCollectCacheHandler(c config.Config) *MsgUpdateUserCollectCacheHandler {
	logger := log.GetSugaredLogger()
	rdb, err := apollo.GetRedisClient("user.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	userDB, err := apollo.GetMysqlDB("user.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
	}

	return &MsgUpdateUserCollectCacheHandler{
		Rdb: rdb,

		UserModel: query.Use(userDB),
	}
}

func NewScheduleUpdateUserCollectRecordHandler(c config.Config) *ScheduleUpdateUserCollectRecordHandler {
	logger := log.GetSugaredLogger()
	rdb, err := apollo.GetRedisClient("user.yaml")
	if err != nil {
		logger.Fatalf("initialize redis failed, err: %v", err)
	}

	userDB, err := apollo.GetMysqlDB("user.yaml")
	if err != nil {
		logger.Fatalf("initialize user mysql failed, err: %v", err)
	}

	return &ScheduleUpdateUserCollectRecordHandler{
		Rdb: rdb,

		UserModel: query.Use(userDB),
	}
}

func (l *MsgCreateUserSubjectHandler) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
	var payload job.MsgCreateUserSubjectPayload
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal [MsgCreateUserSubjectPayload] failed, err: %v", err)
	}

	userSubjectModel := l.UserModel.UserSubject

	userSubject, err := userSubjectModel.WithContext(ctx).
		Where(userSubjectModel.Username.Eq(payload.Username),
			userSubjectModel.Nickname.Eq(payload.Nickname),
			userSubjectModel.Password.Eq(payload.Password),
			userSubjectModel.CreateTime.Eq(payload.CreateTime),
			userSubjectModel.UpdateTime.Eq(payload.UpdateTime)).
		FirstOrCreate()
	if err != nil {
		return fmt.Errorf("update [user_subject] record failed, err: %v", err)
	}

	userSubjectProto := &pb.UserSubject{
		Id:         userSubject.ID,
		Username:   userSubject.Username,
		Password:   userSubject.Password,
		Nickname:   userSubject.Nickname,
		CreateTime: userSubject.CreateTime.String(),
		UpdateTime: userSubject.UpdateTime.String(),
	}

	userSubjectBytes, err := proto.Marshal(userSubjectProto)
	if err != nil {
		return fmt.Errorf("marshal [userSubjectProto] into proto failed, err; %v", err)
	}

	err = l.Rdb.Set(ctx,
		fmt.Sprintf("user_subject_%d", userSubject.ID),
		userSubjectBytes,
		time.Second*86400).Err()
	if err != nil {
		return fmt.Errorf("update [user_subject] cache failed, err: %v", err)
	}

	err = l.Rdb.Set(ctx,
		fmt.Sprintf("user_login_%s", userSubject.Username),
		fmt.Sprintf("%d:%s", userSubject.ID, userSubject.Password),
		time.Second*86400).Err()
	if err != nil {
		return fmt.Errorf("update [user_login] cache failed, err: %v", err)
	}

	return nil
}

func (l *MsgUpdateUserSubjectRecordHandler) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
	var payload job.UserSubjectPayload
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal [UserSubjectPayload] failed, err: %v", err)
	}

	userSubjectModel := l.UserModel.UserSubject

	userSubject := &model.UserSubject{}

	err = structx.SyncWithNoZero(payload, userSubject)
	if err != nil {
		return fmt.Errorf("sync struct [userSubject] failed, err: %v", err)
	}

	_, err = userSubjectModel.WithContext(ctx).
		Where(userSubjectModel.ID.Eq(payload.Id)).
		Updates(userSubject)
	if err != nil {
		return fmt.Errorf("update [user_subject] record failed, err: %v", err)
	}

	return nil
}

func (l *MsgUpdateUserSubjectCacheHandler) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
	var payload job.UserSubjectPayload
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal [UserSubjectPayload] failed, err: %v", err)
	}

	userSubjectProto := &pb.UserSubject{}

	userSubjectBytes, err := l.Rdb.Get(ctx,
		fmt.Sprintf("user_subject_%d", payload.Id)).Bytes()
	switch err {
	case redis.Nil:

	case nil:
		err = proto.Unmarshal(userSubjectBytes, userSubjectProto)
		if err != nil {
			return fmt.Errorf("unmarshal [userSubjectProto] failed, err: %v", err)
		}

	default:
		if err != nil {
			return fmt.Errorf("get [user_subject] cache failed, err: %v", err)
		}
	}

	err = structx.SyncWithNoZero(payload, userSubjectProto)
	if err != nil {
		return fmt.Errorf("sync struct [userSubjectProto] failed, err: %v", err)
	}

	userSubjectBytes, err = proto.Marshal(userSubjectProto)
	if err != nil {
		return fmt.Errorf("marshal [userSubjectProto] failed, err: %v", err)
	}

	err = l.Rdb.Set(ctx,
		fmt.Sprintf("user_subject_%d", userSubjectProto.Id),
		userSubjectBytes,
		time.Second*86400).Err()
	if err != nil {
		return fmt.Errorf("update [user_subject] cache failed, err: %v", err)
	}

	return nil
}

func (l *MsgAddUserSubjectCacheHandler) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
	var payload job.MsgAddUserSubjectCachePayload
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal MsgAddUserSubjectCachePayload failed, err: %v", err)
	}

	userSubjectProto := &pb.UserSubject{}

	// TODO: 待使用分布式锁, 防止脏读
	userSubjectBytes, err := l.Rdb.Get(ctx,
		fmt.Sprintf("user_subject_%d", payload.Id)).Bytes()
	switch err {
	case redis.Nil:
		userSubjectModel := l.UserModel.UserSubject

		userSubject, err := userSubjectModel.WithContext(ctx).
			Select(userSubjectModel.ID, userSubjectModel.Vip, userSubjectModel.Follower).
			Where(userSubjectModel.ID.Eq(payload.Id)).
			First()
		if err != nil {
			return fmt.Errorf("query [user_subject] record failed, err: %v", err)
		}

		err = structx.SyncWithNoZero(*userSubject, userSubjectProto)
		if err != nil {
			return fmt.Errorf("sync struct [pb.UserSubject] failead, err: %v", err)
		}

	case nil:
		err = proto.Unmarshal(userSubjectBytes, userSubjectProto)
		if err != nil {
			return fmt.Errorf("unmarshal [userSubjectProto] failed, err: %v", err)
		}

	default:
		return fmt.Errorf("get [user_subject] cache failed, err: %v", err)
	}
	if err != nil {
	}

	if payload.Vip != 0 {
		userSubjectProto.Vip = userSubjectProto.Vip + payload.Vip
	}
	if payload.Follower != 0 {
		userSubjectProto.Follower = userSubjectProto.Follower + payload.Follower
	}

	userSubjectBytes, err = proto.Marshal(userSubjectProto)
	if err != nil {
		return fmt.Errorf("marshal [userSubjectProto] failed, err: %v", err)
	}

	err = l.Rdb.Set(ctx,
		fmt.Sprintf("user_subject_%d", userSubjectProto.Id),
		userSubjectBytes,
		time.Second*86400).Err()
	if err != nil {
		return fmt.Errorf("update [user_subject] cache failed, err: %v", err)
	}

	return nil
}

func (l *ScheduleUpdateUserSubjectRecordHandler) ProcessTask(ctx context.Context, _ *asynq.Task) (err error) {
	members, err := l.Rdb.SMembers(ctx,
		"user_follower").Result()
	if err != nil {
		return fmt.Errorf("get [user_follower] member failed, err: %v", err)
	}
	l.Rdb.Del(ctx,
		fmt.Sprintf("user_follower"))

	userSubjectModel := l.UserModel.UserSubject
	for _, member := range members {
		followerCount, err := l.Rdb.Get(ctx,
			fmt.Sprintf("user_follower_%s", member)).Int()
		if err != nil {
			return fmt.Errorf("get [user_follower] cnt failed, err: %v", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("user_follwer_%s", member)).Err()
		if err != nil {
			return fmt.Errorf("del [user_follower] cnt failed, err: %v", err)
		}

		userSubject, err := userSubjectModel.WithContext(ctx).
			Select(userSubjectModel.ID, userSubjectModel.Follower).
			Where(userSubjectModel.ID.Eq(cast.ToInt64(member))).
			First()
		if err != nil {
			return fmt.Errorf("get [user_subject] record failed, err: %v", err)
		}

		_, err = userSubjectModel.WithContext(ctx).
			Select(userSubjectModel.ID, userSubjectModel.Follower).
			Where(userSubjectModel.ID.Eq(cast.ToInt64(member))).
			Update(userSubjectModel.Follower, int(userSubject.Follower)+followerCount)
		if err != nil {
			return fmt.Errorf("update [user_subject] record failed, err: %v", err)
		}
	}

	return nil
}

func (l *MsgUpdateUserCollectCacheHandler) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
	return nil
}
