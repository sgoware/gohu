package logic

import (
	"context"
	"main/app/service/notification/rpc/info/info"

	"main/app/service/notification/api/internal/svc"
	"main/app/service/notification/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotificationLogic {
	return &GetNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNotificationLogic) GetNotification(req *types.GetNotificationReq) (resp *types.GetNotificationRes, err error) {
	res, _ := l.svcCtx.InfoRpcClient.GetNotification(l.ctx, &info.GetNotificationReq{NotificationId: req.NotificationId})
	return &types.GetNotificationRes{
		Code: res.Code,
		Msg:  res.Msg,
		Ok:   res.Ok,
		Data: types.GetNotificationResData{
			NotificationSubject: types.NotificationSubject{
				Id:          res.Data.NotificationSubject.Id,
				UserId:      res.Data.NotificationSubject.UserId,
				MessageType: res.Data.NotificationSubject.MessageType,
				CreateTime:  res.Data.NotificationSubject.CreateTime,
				UpdateTime:  res.Data.NotificationSubject.UpdateTime,
			},
			NotificationContent: types.NotificationContent{
				SubjectId:  res.Data.NotificationContent.SubjectId,
				Title:      res.Data.NotificationContent.Title,
				Content:    res.Data.NotificationContent.Content,
				Url:        res.Data.NotificationContent.Url,
				Meta:       res.Data.NotificationContent.Meta,
				Attrs:      res.Data.NotificationContent.Attrs,
				CreateTime: res.Data.NotificationContent.CreateTime,
				UpdateTime: res.Data.NotificationContent.UpdateTime,
			},
		},
	}, nil
}
