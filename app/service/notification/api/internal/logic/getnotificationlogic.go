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
		Data: res.Data.String(),
	}, nil
}
