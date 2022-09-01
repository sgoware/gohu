package logic

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"main/app/common/log"
	"net/http"
	"time"

	"main/app/service/notification/rpc/info/internal/svc"
	"main/app/service/notification/rpc/info/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotificationLogic {
	return &GetNotificationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNotificationLogic) GetNotification(in *pb.GetNotificationReq) (res *pb.GetNotificationRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	resData := &pb.GetNotificationRes_Data{}

	notificationSubjectBytes, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("notificationSubject_%d", in.NotificationId)).Bytes()
	if err != nil {
		logger.Errorf("get notificationSubject cache failed, err: %v", err)

		notificationSubjectModel := l.svcCtx.NotificationModel.NotificationSubject

		notificationSubject, err := notificationSubjectModel.WithContext(l.ctx).
			Where(notificationSubjectModel.ID.Eq(in.NotificationId)).
			First()
		if err != nil {
			logger.Errorf("get notificationSubject failed in mysql failed, err: %v", err)
			res = &pb.GetNotificationRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		notificationSubjectProto := &pb.NotificationSubject{
			Id:          notificationSubject.ID,
			UserId:      notificationSubject.UserID,
			MessageType: notificationSubject.MessageType,
			CreateTime:  notificationSubject.CreateTime.String(),
			UpdateTime:  notificationSubject.UpdateTime.String(),
		}

		resData.NotificationSubject = notificationSubjectProto

		notificationSubjectBytes, err = proto.Marshal(notificationSubjectProto)
		if err != nil {
			logger.Errorf("marshal notificationSubjectProto failed, err: %v", err)
		} else {
			l.svcCtx.Rdb.Set(l.ctx,
				fmt.Sprintf("notificationSubject_%d", notificationSubject.ID),
				notificationSubjectBytes,
				time.Second*86400)
		}
	} else {
		notificationSubjectProto := &pb.NotificationSubject{}
		err = proto.Unmarshal(notificationSubjectBytes, notificationSubjectProto)
		if err != nil {
			logger.Errorf("unmarshal notificationSubjectBytes failed, err: %v", err)
		}

		resData.NotificationSubject = notificationSubjectProto
	}

	notificationContentBytes, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("notificationContent_%d", in.NotificationId)).Bytes()
	if err != nil {
		logger.Errorf("get notificationContent cache failed, err: %v")

		notificationContentModel := l.svcCtx.NotificationModel.NotificationContent

		notificationContent, err := notificationContentModel.WithContext(l.ctx).
			Where(notificationContentModel.SubjectID.Eq(in.NotificationId)).
			First()
		if err != nil {
			logger.Errorf("get notificationContent in mysql failed, err: %v", err)
			res = &pb.GetNotificationRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		notificationContentProto := &pb.NotificationContent{
			SubjectId:  notificationContent.SubjectID,
			Title:      notificationContent.Title,
			Content:    notificationContent.Content,
			Url:        notificationContent.URL,
			Attrs:      notificationContent.Attrs,
			CreateTime: notificationContent.CreateTime.String(),
			UpdateTime: notificationContent.UpdateTime.String(),
		}

		resData.NotificationContent = notificationContentProto

		notificationContentBytes, err = proto.Marshal(notificationContentProto)
		if err != nil {
			logger.Errorf("marshal notificationContentProto failed, err: %v", err)
		} else {
			l.svcCtx.Rdb.Set(l.ctx,
				fmt.Sprintf("notificationContent_%d", notificationContent.SubjectID),
				notificationContentBytes,
				time.Second*86400)
		}
	} else {
		notificationContentProto := &pb.NotificationContent{}
		err = proto.Unmarshal(notificationContentBytes, notificationContentProto)
		if err != nil {
			logger.Errorf("unmarshal notificationContentProto failed, err: %v", err)
			res = &pb.GetNotificationRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", err)
			return res, nil
		}
	}

	res = &pb.GetNotificationRes{
		Code: http.StatusOK,
		Msg:  "get notification successfully",
		Ok:   true,
		Data: resData,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
