package logic

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"main/app/common/log"
	"main/app/service/notification/rpc/crud/internal/svc"
	"main/app/service/notification/rpc/crud/pb"
	modelpb "main/app/service/question/dao/pb"
	"net/http"
	"time"

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

	notificationSubjectProto := &modelpb.NotificationSubject{
		Id:          notificationSubject.ID,
		UserId:      notificationSubject.UserID,
		MessageType: notificationSubject.MessageType,
		CreateTime:  notificationSubject.CreateTime.String(),
		UpdateTime:  notificationSubject.UpdateTime.String(),
	}
	notificationSubjectBytes, err := proto.Marshal(notificationSubjectProto)
	if err != nil {
		logger.Errorf("marshal proto failed, err: %v", err)
		res = &pb.PublishNotificationRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
	// 设置缓存
	l.svcCtx.Rdb.Set(l.ctx,
		fmt.Sprintf("notification_subject_%d", notificationSubject.ID),
		notificationSubjectBytes,
		time.Second*86400)

	// 设置用户所有的notification_id
	l.svcCtx.Rdb.SAdd(l.ctx,
		fmt.Sprintf("notification_%d_0", notificationSubject.UserID))
	l.svcCtx.Rdb.SAdd(l.ctx,
		fmt.Sprintf("notification_%d_%d", notificationSubject.UserID, notificationSubject.MessageType))

	notificationContent, err := notificationContentModel.WithContext(l.ctx).
		Where(notificationContentModel.SubjectID.Eq(notificationSubject.ID),
			notificationContentModel.Title.Eq(in.Title),
			notificationContentModel.Content.Eq(in.Content),
			notificationContentModel.URL.Eq(in.Url)).
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
	notificationContentProto := &modelpb.NotificationContent{
		SubjectId:  notificationContent.SubjectID,
		Title:      notificationContent.Title,
		Content:    notificationContent.Content,
		Url:        notificationContent.URL,
		Attrs:      notificationContent.Attrs,
		CreateTime: notificationContent.CreateTime.String(),
		UpdateTime: notificationContent.UpdateTime.String(),
	}
	notificationContentBytes, err := proto.Marshal(notificationContentProto)
	if err != nil {
		logger.Errorf("marshal proto failed, err: %v", err)
		res = &pb.PublishNotificationRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
	l.svcCtx.Rdb.Set(l.ctx,
		fmt.Sprintf("notificationContent_%d", notificationContent.SubjectID),
		notificationContentBytes,
		time.Second*86400)

	res = &pb.PublishNotificationRes{
		Code: http.StatusOK,
		Msg:  "publish notification successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
