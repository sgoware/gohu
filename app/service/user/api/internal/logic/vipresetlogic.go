package logic

import (
	"context"
	"github.com/spf13/cast"
	"main/app/service/user/rpc/vip/vip"

	"main/app/service/user/api/internal/svc"
	"main/app/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VipResetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVipResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VipResetLogic {
	return &VipResetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VipResetLogic) VipReset(req *types.VipResetReq) (*types.VipResetRes, error) {
	userId := l.ctx.Value("user_id")
	res, _ := l.svcCtx.VipRpcClient.Reset(l.ctx, &vip.ResetReq{Uid: cast.ToString(userId)})

	return &types.VipResetRes{
		Code: int(res.Code),
		Msg:  res.Msg,
		Ok:   res.Ok,
	}, nil
}
