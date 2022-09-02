package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"main/app/common/log"
	"main/app/service/mq/asynq/processor/job"
	modelpb "main/app/service/user/dao/pb"
	"main/app/service/user/rpc/info/internal/svc"
	"main/app/service/user/rpc/info/pb"
	"main/app/utils/structx"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPersonalInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPersonalInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPersonalInfoLogic {
	return &GetPersonalInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPersonalInfoLogic) GetPersonalInfo(in *pb.GetPersonalInfoReq) (res *pb.GetPersonalInfoRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	userSubjectBytes, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("user_subject_%d", in.UserId)).Bytes()
	if err == nil {
		// 查找缓存成功
		rpcResData := &pb.GetPersonalInfoRes_Data{}

		userSubjectProto := &modelpb.UserSubject{}

		err = proto.Unmarshal(userSubjectBytes, userSubjectProto)
		if err != nil {
			logger.Errorf("unmarshal [rpcResData] failed, err: %v", err)
		} else {
			err = structx.SyncWithNoZero(*userSubjectProto, rpcResData)
			if err != nil {
				logger.Errorf("sync struct [GetPersonalInfoRes_Data] failed, err: %v", err)
				res = &pb.GetPersonalInfoRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %v", err)
				return res, nil
			}
			res = &pb.GetPersonalInfoRes{
				Code: http.StatusOK,
				Msg:  "get personal info successfully",
				Ok:   true,
				Data: rpcResData,
			}
			logger.Debugf("send message: %v", err)
			return res, nil
		}
	}
	logger.Errorf("get [user_subject] cache failed, err: %v", err)

	// 在数据库中查找
	userSubjectModel := l.svcCtx.UserModel.UserSubject

	userSubject, err := userSubjectModel.WithContext(l.ctx).
		Where(userSubjectModel.ID.Eq(in.UserId)).
		First()
	switch err {
	case gorm.ErrRecordNotFound:
		res = &pb.GetPersonalInfoRes{
			Code: http.StatusForbidden,
			Msg:  "user not found",
			Ok:   false,
			Data: nil,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil

	case nil:

	default:
		logger.Debugf("get personal info failed, err: mysql err, %v", err)
		res = &pb.GetPersonalInfoRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
			Data: nil,
		}
		return res, nil
	}

	// 更新缓存
	payload := &job.UserSubjectPayload{}
	err = structx.SyncWithNoZero(*userSubject, payload)
	if err != nil {
		logger.Errorf("sync struct [UserSubjectPayload] failed, err: %v", err)
	}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		logger.Errorf("marshal [payload] to json failed, err: %v", err)
	}

	_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(job.MsgUpdateUserSubjectCacheTask, payloadJson))
	if err != nil {
		logger.Errorf("create [MsgUpdateUserSubjectCacheTask] insert queue failed, err: %v", err)
	}

	res = &pb.GetPersonalInfoRes{
		Code: http.StatusOK,
		Msg:  "get personal info successfully",
		Ok:   true,
		Data: &pb.GetPersonalInfoRes_Data{
			Username:   userSubject.Username,
			Nickname:   userSubject.Nickname,
			Email:      userSubject.Email,
			Phone:      userSubject.Phone,
			LastIp:     userSubject.LastIP,
			Vip:        userSubject.Vip,
			Follower:   userSubject.Follower,
			State:      userSubject.State,
			CreateTime: userSubject.CreateTime.String(),
			UpdateTime: userSubject.UpdateTime.String(),
		},
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
