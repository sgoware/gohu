package user

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
	"main/app/service/user/dao/model"
	"main/app/service/user/dao/pb"
	"main/app/service/user/dao/query"
	"main/app/utils/structx"
	"strings"
	"time"
)

type MsgCreateUserSubjectHandler struct {
	Rdb         *redis.Client
	UserModel   *query.Query
	IdGenerator *idgen.DefaultIdGenerator
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
	Rdb         *redis.Client
	UserModel   *query.Query
	IdGenerator *idgen.DefaultIdGenerator
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

	idGenerator, err := apollo.NewIdGenerator("user.yaml")
	if err != nil {
		logger.Fatalf("get idGenerator failed, err: %v", err)
	}

	return &MsgCreateUserSubjectHandler{
		Rdb: rdb,

		UserModel: query.Use(userDB),

		IdGenerator: idGenerator,
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

	idGenerator, err := apollo.NewIdGenerator("user.yaml")
	if err != nil {
		logger.Errorf("get idGenerator failed, err: %v", err)
	}

	return &ScheduleUpdateUserCollectRecordHandler{
		Rdb: rdb,

		UserModel: query.Use(userDB),

		IdGenerator: idGenerator,
	}
}

func (l *MsgCreateUserSubjectHandler) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
	var payload job.MsgCreateUserSubjectPayload
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal [MsgCreateUserSubjectPayload] failed, err: %v", err)
	}

	userSubjectId := l.IdGenerator.NewLong()

	userSubjectModel := l.UserModel.UserSubject

	now := time.Now()

	err = userSubjectModel.WithContext(ctx).
		Create(&model.UserSubject{
			ID:         userSubjectId,
			Username:   payload.Username,
			Password:   payload.Password,
			Nickname:   payload.Nickname,
			CreateTime: now,
			UpdateTime: now,
		})
	if err != nil {
		return fmt.Errorf("create [user_subject] record failed, err: %v", err)
	}

	userSubjectProto := &pb.UserSubject{
		Id:         userSubjectId,
		Username:   payload.Username,
		Password:   payload.Password,
		Nickname:   payload.Nickname,
		CreateTime: now.String(),
		UpdateTime: now.String(),
	}

	userSubjectBytes, err := proto.Marshal(userSubjectProto)
	if err != nil {
		return fmt.Errorf("marshal [userSubjectProto] into proto failed, err: %v", err)
	}

	err = l.Rdb.Set(ctx,
		fmt.Sprintf("user_subject_%d", userSubjectId),
		userSubjectBytes,
		time.Second*86400).Err()
	if err != nil {
		return fmt.Errorf("update [user_subject] cache failed, err: %v", err)
	}

	err = l.Rdb.Set(ctx,
		fmt.Sprintf("user_login_%d", userSubjectId),
		fmt.Sprintf("%d:%s", userSubjectId, payload.Password),
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
		"user_follower_cnt_set").Result()
	if err != nil {
		return fmt.Errorf("get [user_follower] member failed, err: %v", err)
	}
	l.Rdb.Del(ctx,
		fmt.Sprintf("user_follower_cnt_set"))

	userSubjectModel := l.UserModel.UserSubject
	for _, member := range members {
		followerCount, err := l.Rdb.Get(ctx,
			fmt.Sprintf("user_follower_cnt_%s", member)).Int()
		if err != nil {
			return fmt.Errorf("get [user_follower] cnt failed, err: %v", err)
		}

		err = l.Rdb.Del(ctx,
			fmt.Sprintf("user_follower_cnt_%s", member)).Err()
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

func (l *ScheduleUpdateUserCollectRecordHandler) ProcessTask(ctx context.Context, _ *asynq.Task) (err error) {
	userCollectionModel := l.UserModel.UserCollection
	for {
		cmd, err := l.Rdb.RPop(ctx, "user_collection_list").Result()
		switch err {
		case redis.Nil:
			return nil
		case nil:

		default:
			return fmt.Errorf("get [user_collect] list member failed, err: %v", err)
		}
		output := strings.Split(cmd, "_")
		userId := cast.ToInt64(output[1])
		collectType := cast.ToInt32(output[2])
		objType := cast.ToInt32(output[3])
		objId := cast.ToInt64(output[4])
		if output[0] == "0" {
			// 创建操作

			userCollectionId := l.IdGenerator.NewLong()

			err = userCollectionModel.WithContext(ctx).
				Create(&model.UserCollection{
					ID:          userCollectionId,
					UserID:      userId,
					CollectType: collectType,
					ObjType:     objType,
					ObjID:       objId,
				})
			if err != nil {
				return fmt.Errorf("create [user_collect] record failed, err: %v", err)
			}
		} else {
			// 删除操作
			_, err = userCollectionModel.WithContext(ctx).
				Where(userCollectionModel.UserID.Eq(userId),
					userCollectionModel.CollectType.Eq(collectType),
					userCollectionModel.ObjType.Eq(objType),
					userCollectionModel.ObjID.Eq(objId)).
				Delete()
			if err != nil {
				return fmt.Errorf("delete [user_collect] record failed, err: %v", err)
			}
		}
	}
}
