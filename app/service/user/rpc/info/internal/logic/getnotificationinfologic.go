package logic

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"main/app/common/log"
	"main/app/service/notification/dao/model"
	"net/http"

	"main/app/service/user/rpc/info/internal/svc"
	"main/app/service/user/rpc/info/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNotificationInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNotificationInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotificationInfoLogic {
	return &GetNotificationInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNotificationInfoLogic) GetNotificationInfo(in *pb.GetNotificationInfoReq) (res *pb.GetNotificationInfoRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	var notificationIds []int64

	notificationIdsCache, err := l.svcCtx.Rdb.SMembers(l.ctx,
		fmt.Sprintf("notification_%d_%d", in.UserId, in.MessageType)).Result()
	if err == nil {
		for _, notificationIdCache := range notificationIdsCache {
			notificationIds = append(notificationIds, cast.ToInt64(notificationIdCache))
		}
	} else {
		logger.Errorf("get notification ids cache %d failed, err: %v", in.UserId, err)

		notificationSubjectModel := l.svcCtx.NotificationModel.NotificationSubject

		notificationSubjects := make([]*model.NotificationSubject, 0)
		if in.MessageType == 0 {
			notificationSubjects, err = notificationSubjectModel.WithContext(l.ctx).
				Select(notificationSubjectModel.UserID, notificationSubjectModel.ID).
				Where(notificationSubjectModel.UserID.Eq(in.UserId)).
				Find()
			if err != nil {
				logger.Errorf("get notification ids in mysql failed, err: %v", err)
				res = &pb.GetNotificationInfoRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %v", res.String())
				return res, nil
			}
		} else {
			notificationSubjects, err = notificationSubjectModel.WithContext(l.ctx).
				Select(notificationSubjectModel.UserID, notificationSubjectModel.MessageType, notificationSubjectModel.ID).
				Where(notificationSubjectModel.UserID.Eq(in.UserId),
					notificationSubjectModel.MessageType.Eq(in.MessageType)).
				Find()
			if err != nil {
				logger.Errorf("get notification ids in mysql failed, err: %v", err)
				res = &pb.GetNotificationInfoRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %v", res.String())
				return res, nil
			}
		}
		for _, notificationSubject := range notificationSubjects {
			notificationIds = append(notificationIds, notificationSubject.ID)
			l.svcCtx.Rdb.SAdd(l.ctx,
				fmt.Sprintf("notification_%d_0", in.UserId),
				notificationSubject.ID)
			l.svcCtx.Rdb.SAdd(l.ctx,
				fmt.Sprintf("notification_%d_%d", in.UserId, in.MessageType),
				notificationSubject.ID)
		}
	}

	res = &pb.GetNotificationInfoRes{
		Code: http.StatusOK,
		Msg:  "get notification ids successfully",
		Ok:   true,
		Data: &pb.GetNotificationInfoRes_Data{NotificationIds: notificationIds},
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
