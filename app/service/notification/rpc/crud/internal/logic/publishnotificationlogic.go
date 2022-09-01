package logic

import (
	"context"
	"main/app/common/log"
	"main/app/service/notification/dao/model"
	"main/app/service/notification/rpc/crud/internal/svc"
	"main/app/service/notification/rpc/crud/pb"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishNotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishNotificationLogic {
	return &PublishNotificationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishNotificationLogic) PublishNotification(in *pb.PublishNotificationReq) (res *pb.PublishNotificationRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	notificationSubjectModel := l.svcCtx.NotificationModel.NotificationSubject
	notificationContentModel := l.svcCtx.NotificationModel.NotificationContent

	notificationSubject, err := notificationSubjectModel.WithContext(l.ctx).
		Where(notificationSubjectModel.UserID.Eq(in.UserId),
			notificationSubjectModel.MessageType.Eq(in.MessageType)).
		FirstOrCreate()
	if err != nil {
		logger.Errorf("publish notification failed, err: mysql err, %v", err)
		res = &pb.PublishNotificationRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	err = notificationContentModel.WithContext(l.ctx).
		Create(&model.NotificationContent{
			SubjectID: notificationSubject.ID,
			Title:     in.Title,
			Content:   in.Content,
			URL:       in.Url,
		})
	if err != nil {
		logger.Errorf("publish notification failed, err: mysql err, %v", err)
		res = &pb.PublishNotificationRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	res = &pb.PublishNotificationRes{
		Code: http.StatusOK,
		Msg:  "publish notification successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
