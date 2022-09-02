package logic

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/gorm"
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
		if len(notificationIdsCache) > 1 {
			for _, notificationIdCache := range notificationIdsCache {
				if notificationIdCache != "0" {
					notificationIds = append(notificationIds, cast.ToInt64(notificationIdCache))
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
	} else {
		logger.Errorf("get notification ids cache %d failed, err: %v", in.UserId, err)
	}

	notificationSubjectModel := l.svcCtx.NotificationModel.NotificationSubject

	notificationSubjects := make([]*model.NotificationSubject, 0)
	if in.MessageType == 0 {
		notificationSubjects, err = notificationSubjectModel.WithContext(l.ctx).
			Select(notificationSubjectModel.UserID, notificationSubjectModel.ID).
			Where(notificationSubjectModel.UserID.Eq(in.UserId)).
			Find()
		switch err {
		case gorm.ErrRecordNotFound:
			// 设置空缓存
			l.svcCtx.Rdb.SAdd(l.ctx,
				fmt.Sprintf("notification_%d_0", in.UserId),
				0)
			res = &pb.GetNotificationInfoRes{
				Code: http.StatusOK,
				Msg:  "get notification ids successfully",
				Ok:   true,
				Data: &pb.GetNotificationInfoRes_Data{NotificationIds: nil},
			}
			logger.Debugf("send message: %v", err)
			return res, nil

		case nil:

		default:
			logger.Errorf("query [notification_subject] record failed ,err: %v", err)
			res = &pb.GetNotificationInfoRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
			}
			logger.Debugf("send message: %v", err)
			return res, nil
		}
	} else {
		notificationSubjects, err = notificationSubjectModel.WithContext(l.ctx).
			Select(notificationSubjectModel.UserID, notificationSubjectModel.MessageType, notificationSubjectModel.ID).
			Where(notificationSubjectModel.UserID.Eq(in.UserId),
				notificationSubjectModel.MessageType.Eq(in.MessageType)).
			Find()
		switch err {
		case gorm.ErrRecordNotFound:
			l.svcCtx.Rdb.SAdd(l.ctx,
				fmt.Sprintf("notification_%d_%d", in.UserId, in.MessageType),
				0)
			res = &pb.GetNotificationInfoRes{
				Code: http.StatusOK,
				Msg:  "get notification ids successfully",
				Ok:   true,
				Data: &pb.GetNotificationInfoRes_Data{NotificationIds: nil},
			}
			logger.Debugf("send message: %v", err)
			return res, nil

		case nil:

		default:
			logger.Errorf("query [notification_subject] record failed ,err: %v", err)
			res = &pb.GetNotificationInfoRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
			}
			logger.Debugf("send message: %v", err)
			return res, nil
		}
	}

	for _, notificationSubject := range notificationSubjects {
		notificationIds = append(notificationIds, notificationSubject.ID)
		l.svcCtx.Rdb.SAdd(l.ctx,
			fmt.Sprintf("notification_%d_0", in.UserId),
			notificationSubject.ID)
		l.svcCtx.Rdb.SAdd(l.ctx,
			fmt.Sprintf("notification_%d_%d", in.UserId, notificationSubject.MessageType),
			notificationSubject.ID)
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
