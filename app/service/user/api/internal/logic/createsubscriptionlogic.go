package logic

import (
	"context"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"main/app/service/user/rpc/crud/crud"

	"main/app/service/user/api/internal/svc"
	"main/app/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSubscriptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSubscriptionLogic {
	return &CreateSubscriptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSubscriptionLogic) CreateSubscription(req *types.CreateSubscriptionReq) (resp *types.CreateSubscriptionRes, err error) {
	j := gjson.Parse(cast.ToString(l.ctx.Value("user_details")))
	userId := j.Get("user_id").Int()
	res, _ := l.svcCtx.CrudRpcClient.CreateSubscription(l.ctx, &crud.CreateSubscriptionReq{
		UserId:  userId,
		ObjType: req.ObjType,
		ObjId:   req.ObjId,
	})

	return &types.CreateSubscriptionRes{
		Code: res.Code,
		Msg:  res.Msg,
		Ok:   res.Ok,
	}, nil
}
