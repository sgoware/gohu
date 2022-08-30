package logic

import (
	"context"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
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
	j := gjson.Parse(cast.ToString(l.ctx.Value("user_details")))
	userId := j.Get("user_id").Int()
	res, _ := l.svcCtx.VipRpcClient.Upgrade(l.ctx, &vip.UpgradeReq{Id: userId})

	return &types.VipUpgradeRes{
		Code: res.Code,
		Msg:  res.Msg,
		Ok:   res.Ok,
		Data: types.VipUpgradeResData{VipLevel: res.Data.VipLevel},
	}, nil
}
