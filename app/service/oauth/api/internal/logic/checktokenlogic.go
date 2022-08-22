package logic

import (
	"context"

	"gohu/app/service/oauth/api/internal/svc"
	"gohu/app/service/oauth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckTokenLogic {
	return &CheckTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckTokenLogic) CheckToken(req *types.CheckTokenReq) (resp *types.CheckTokenRes, err error) {
	// todo: add your logic here and delete this line

	return
}
