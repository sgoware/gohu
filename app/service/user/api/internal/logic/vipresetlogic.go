package logic

import (
	"context"

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

func (l *VipResetLogic) VipReset(req *types.VipResetReq) (resp *types.VipResetRes, err error) {
	// todo: add your logic here and delete this line

	return
}
