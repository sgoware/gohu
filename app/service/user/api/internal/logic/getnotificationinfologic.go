package logic

import (
	"context"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"main/app/service/user/rpc/info/info"

	"main/app/service/user/api/internal/svc"
	"main/app/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNotificationInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNotificationInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotificationInfoLogic {
	return &GetNotificationInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNotificationInfoLogic) GetNotificationInfo(req *types.GetNotificationInfoReq) (resp *types.GetNotificationInfoRes, err error) {
	j := gjson.Parse(cast.ToString(l.ctx.Value("user_details")))
	userId := j.Get("user_id").Int()
	res, _ := l.svcCtx.InfoRpcClient.GetNotificationInfo(l.ctx, &info.GetNotificationInfoReq{
		UserId:      userId,
		MessageType: req.MessageType,
	})
	return &types.GetNotificationInfoRes{
		Code: res.Code,
		Msg:  res.Msg,
		Ok:   res.Ok,
		Data: types.GetNotificationInfoResData{
			MessageIds: res.Data.MessageId,
		},
	}, nil
}
