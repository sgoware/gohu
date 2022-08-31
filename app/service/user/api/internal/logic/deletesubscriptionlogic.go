package logic

import (
	"context"
	"main/app/service/user/rpc/crud/crud"

	"main/app/service/user/api/internal/svc"
	"main/app/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSubscriptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSubscriptionLogic {
	return &DeleteSubscriptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSubscriptionLogic) DeleteSubscription(req *types.DeleteSubscriptionReq) (resp *types.DeleteSubscriptionRes, err error) {
	res, _ := l.svcCtx.CrudRpcClient.DeleteSubscription(l.ctx, &crud.DeleteSubscriptionReq{
		SubscriptionId: req.SubscriptionId})
	return &types.DeleteSubscriptionRes{
		Code: res.Code,
		Msg:  res.Msg,
		Ok:   res.Ok,
	}, nil
}
