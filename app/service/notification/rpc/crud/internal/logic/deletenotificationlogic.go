package logic

import (
	"context"
	"main/app/common/log"
	"net/http"

	"main/app/service/notification/rpc/crud/internal/svc"
	"main/app/service/notification/rpc/crud/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNotificationLogic {
	return &DeleteNotificationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteNotificationLogic) DeleteNotification(in *pb.DeleteNotificationReq) (res *pb.DeleteNotificationRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	notificationSubjectModel := l.svcCtx.NotificationModel.NotificationSubject

	_, err = notificationSubjectModel.WithContext(l.ctx).
		Where(notificationSubjectModel.ID.Eq(in.MessageId)).
		Delete()
	if err != nil {
		logger.Errorf("delete notification failed, err: mysql err, %v", err)
		res = &pb.DeleteNotificationRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	res = &pb.DeleteNotificationRes{
		Code: http.StatusOK,
		Msg:  "delete notification successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
