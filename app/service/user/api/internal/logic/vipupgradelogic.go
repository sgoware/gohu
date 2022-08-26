package logic

import (
	"context"
	"github.com/spf13/cast"
	"main/app/service/user/rpc/vip/vip"

	"main/app/service/user/api/internal/svc"
	"main/app/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VipUpgradeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVipUpgradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VipUpgradeLogic {
	return &VipUpgradeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VipUpgradeLogic) VipUpgrade(req *types.VipUpgradeReq) (*types.VipUpgradeRes, error) {
	userId := l.ctx.Value("user_id")
	res, _ := l.svcCtx.VipRpcClient.Upgrade(l.ctx, &vip.UpgradeReq{Id: cast.ToInt64(userId)})

	return &types.VipUpgradeRes{
		Code: int(res.Code),
		Msg:  res.Msg,
		Data: types.VipUpgradeResData{VipLevel: int(res.Data.VipLevel)},
	}, nil
}
