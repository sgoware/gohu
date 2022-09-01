package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"google.golang.org/protobuf/proto"
	"main/app/service/mq/asynq/processor/internal/svc"
	"main/app/service/mq/asynq/processor/job"
	"main/app/service/user/dao/model"
	"main/app/service/user/dao/pb"
	"main/app/utils/structx"
	"time"
)

type MsgCreateUserSubjectHandler struct {
	svcCtx *svc.ServiceContext
}

type MsgUpdateUserSubjectRecordHandler struct {
	svcCtx *svc.ServiceContext
}

type MsgUpdateUserSubjectCacheHandler struct {
	svcCtx *svc.ServiceContext
}

type MsgAddUserSubjectCacheHandler struct {
	svcCtx *svc.ServiceContext
}

type ScheduleUpdateUserSubjectRecordHandler struct {
	svcCtx *svc.ServiceContext
}

func NewCreateUserSubjectRecordHandler(svcCtx *svc.ServiceContext) *MsgCreateUserSubjectHandler {
	return &MsgCreateUserSubjectHandler{
		svcCtx: svcCtx,
	}
}

func NewUpdateUserSubjectRecordHandler(svcCtx *svc.ServiceContext) *MsgUpdateUserSubjectRecordHandler {
	return &MsgUpdateUserSubjectRecordHandler{
		svcCtx: svcCtx,
	}
}

func NewUpdateUserSubjectCacheHandler(svcCtx *svc.ServiceContext) *MsgUpdateUserSubjectCacheHandler {
	return &MsgUpdateUserSubjectCacheHandler{
		svcCtx: svcCtx,
	}
}

func NewMsgAddUserSubjectCacheHandler(svcCtx *svc.ServiceContext) *MsgAddUserSubjectCacheHandler {
	return &MsgAddUserSubjectCacheHandler{
		svcCtx: svcCtx,
	}
}

func NewScheduleUpdateUserSubjectRecordHandler(svcCtx *svc.ServiceContext) *ScheduleUpdateUserSubjectRecordHandler {
	return &ScheduleUpdateUserSubjectRecordHandler{
		svcCtx: svcCtx,
	}
}

func (l *MsgCreateUserSubjectHandler) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
	var payload job.MsgCreateUserSubjectPayload
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal [MsgCreateUserSubjectPayload] failed, err: %v", err)
	}

	userSubjectModel := l.svcCtx.UserModel.UserSubject

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

	err = l.svcCtx.Rdb.Set(ctx,
		fmt.Sprintf("user_subject_%d", userSubject.ID),
		userSubjectBytes,
		time.Second*86400).Err()
	if err != nil {
		return fmt.Errorf("update [user_subject] cache failed, err: %v", err)
	}

	err = l.svcCtx.Rdb.Set(ctx,
		fmt.Sprintf("user_login_%s", userSubject.Username),
		fmt.Sprintf("%d:%s", userSubject.ID, userSubject.Password),
		time.Second*86400).Err()
	if err != nil {
		return fmt.Errorf("update [user_login] cache failed, err: %v", err)
	}

	return nil
}

func (l *MsgUpdateUserSubjectRecordHandler) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
	var payload job.MsgUpdateUserSubjectRecordPayload
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal [MsgUpdateUserSubjectRecordPayload] failed, err: %v", err)
	}

	userSubjectModel := l.svcCtx.UserModel.UserSubject

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
	var payload job.MsgUpdateUserSubjectRecordPayload
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal [MsgUpdateUserSubjectRecordPayload] failed, err: %v", err)
	}

	userSubjectBytes, err := l.svcCtx.Rdb.Get(ctx,
		fmt.Sprintf("user_subject_%d", payload.Id)).Bytes()
	if err != nil {
		return fmt.Errorf("get [user_subject] cache failed, err: %v", err)
	}

	userSubjectProto := &pb.UserSubject{}

	err = proto.Unmarshal(userSubjectBytes, userSubjectProto)
	if err != nil {
		return fmt.Errorf("unmarshal [userSubjectProto] failed, err: %v", err)
	}

	err = structx.SyncWithNoZero(payload, userSubjectProto)
	if err != nil {
		return fmt.Errorf("sync struct [userSubjectProto] failed, err: %v", err)
	}

	userSubjectBytes, err = proto.Marshal(userSubjectProto)
	if err != nil {
		return fmt.Errorf("marshal [userSubjectProto] failed, err: %v", err)
	}

	err = l.svcCtx.Rdb.Set(ctx,
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

	// TODO: 待使用分布式锁, 防止脏读
	userSubjectBytes, err := l.svcCtx.Rdb.Get(ctx,
		fmt.Sprintf("user_subject_%d", payload.Id)).Bytes()
	if err != nil {
		return fmt.Errorf("get [user_subject] cache failed, err: %v", err)
	}

	userSubjectProto := &pb.UserSubject{}

	err = proto.Unmarshal(userSubjectBytes, userSubjectProto)
	if err != nil {
		return fmt.Errorf("unmarshal [userSubjectProto] failed, err: %v", err)
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

	err = l.svcCtx.Rdb.Set(ctx,
		fmt.Sprintf("user_subject_%d", userSubjectProto.Id),
		userSubjectBytes,
		time.Second*86400).Err()
	if err != nil {
		return fmt.Errorf("update [user_subject] cache failed, err: %v", err)
	}

	return nil
}

func (l *ScheduleUpdateUserSubjectRecordHandler) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {

	return nil
}
