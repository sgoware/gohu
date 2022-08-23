package logic

import (
	"context"

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

func (l *VipUpgradeLogic) VipUpgrade(req *types.VipUpgradeReq) (resp *types.VipUpgradeRes, err error) {
	// todo: add your logic here and delete this line

	return
}
