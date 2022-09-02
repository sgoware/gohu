package logic

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	apollo "main/app/common/config"
	"main/app/common/log"
	"main/app/service/notification/dao/model"
	modelpb "main/app/service/notification/dao/pb"
	"main/app/service/notification/rpc/crud/internal/svc"
	"main/app/service/notification/rpc/crud/pb"
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

	idGenerator, err := apollo.NewIdGenerator("notification.yaml")
	if err != nil {
		logger.Errorf("get idGenerator failed, err: %v", err)
		res = &pb.PublishNotificationRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", err)
		return res, nil
	}

	notificationSubjectId := idGenerator.NewLong()
	nowTime := time.Now()

	notificationSubjectModel := l.svcCtx.NotificationModel.NotificationSubject
	notificationContentModel := l.svcCtx.NotificationModel.NotificationContent

	err = notificationSubjectModel.WithContext(l.ctx).
		Create(&model.NotificationSubject{
			ID:          notificationSubjectId,
			UserID:      in.UserId,
			MessageType: in.MessageType,
			CreateTime:  nowTime,
			UpdateTime:  nowTime,
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

	notificationSubjectProto := &modelpb.NotificationSubject{
		Id:          notificationSubjectId,
		UserId:      in.UserId,
		MessageType: in.MessageType,
		CreateTime:  nowTime.String(),
		UpdateTime:  nowTime.String(),
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
		fmt.Sprintf("notification_subject_%d", notificationSubjectId),
		notificationSubjectBytes,
		time.Second*86400)

	// 设置用户所有的notification_id
	l.svcCtx.Rdb.SAdd(l.ctx,
		fmt.Sprintf("notification_%d_0", in.UserId),
		notificationSubjectId)
	l.svcCtx.Rdb.SAdd(l.ctx,
		fmt.Sprintf("notification_%d_%d", in.UserId, in.MessageType),
		notificationSubjectId)

	err = notificationContentModel.WithContext(l.ctx).
		Create(&model.NotificationContent{
			SubjectID:  notificationSubjectId,
			Title:      in.Title,
			Content:    in.Content,
			URL:        in.Url,
			Meta:       "",
			Attrs:      0,
			CreateTime: nowTime,
			UpdateTime: nowTime,
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
	notificationContentProto := &modelpb.NotificationContent{
		SubjectId:  notificationSubjectId,
		Title:      in.Title,
		Content:    in.Content,
		Url:        in.Url,
		Meta:       "",
		Attrs:      0,
		CreateTime: nowTime.String(),
		UpdateTime: nowTime.String(),
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
		fmt.Sprintf("notification_content_%d", notificationSubjectId),
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
