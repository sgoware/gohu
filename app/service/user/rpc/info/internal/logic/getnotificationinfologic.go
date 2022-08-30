package logic

import (
	"context"
	"main/app/common/log"
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

	notificationSubjectModel := l.svcCtx.NotificationModel.NotificationSubject

	switch in.MessageType {
	case 0:
		// 获取全部通知
		notificationSubjects, err := notificationSubjectModel.WithContext(l.ctx).
			Where(notificationSubjectModel.UserID.Eq(in.UserId)).Find()
		if err != nil {
			logger.Errorf("get notification info failed, err: mysql err, %v", err)
			res = &pb.GetNotificationInfoRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
				Data: nil,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		res = &pb.GetNotificationInfoRes{
			Code: http.StatusOK,
			Msg:  "get notification successfully",
			Ok:   true,
			Data: &pb.GetNotificationInfoRes_Data{MessageId: make([]int64, 0)},
		}
		for _, notificationSubject := range notificationSubjects {
			res.Data.MessageId = append(res.Data.MessageId, notificationSubject.ID)
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil

	default:
		// 获取指定类型通知
		notificationSubjects, err := notificationSubjectModel.WithContext(l.ctx).
			Where(notificationSubjectModel.UserID.Eq(in.UserId),
				notificationSubjectModel.MessageType.Eq(in.MessageType)).Find()
		if err != nil {
			logger.Errorf("get notification info failed, err: mysql err, %v", err)
			res = &pb.GetNotificationInfoRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
				Data: nil,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		res = &pb.GetNotificationInfoRes{
			Code: http.StatusOK,
			Msg:  "get notification successfully",
			Ok:   true,
			Data: &pb.GetNotificationInfoRes_Data{MessageId: make([]int64, 0)},
		}
		for _, notificationSubject := range notificationSubjects {
			res.Data.MessageId = append(res.Data.MessageId, notificationSubject.ID)
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
}
