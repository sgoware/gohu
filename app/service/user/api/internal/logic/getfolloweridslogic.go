package logic

import (
	"context"
	"main/app/service/user/rpc/info/info"

	"main/app/service/user/api/internal/svc"
	"main/app/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowerIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFollowerIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowerIdsLogic {
	return &GetFollowerIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFollowerIdsLogic) GetFollowerIds(req *types.GetFollowerIdsReq) (resp *types.GetFollowerIdsRes, err error) {
	rpcRes, _ := l.svcCtx.InfoRpcClient.GetFollower(l.ctx, &info.GetFollowerReq{UserId: req.UserId})

	return &types.GetFollowerIdsRes{
		Code: rpcRes.Code,
		Msg:  rpcRes.Msg,
		Ok:   rpcRes.Ok,
		Data: types.GetFollowerIdsResData{UserIds: rpcRes.Data.UserIds},
	}, nil
}
