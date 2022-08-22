package logic

import (
	"context"

	"gohu/app/service/oauth/api/internal/svc"
	"gohu/app/service/oauth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenLogic {
	return &GetTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenLogic) GetToken(req *types.GetTokenByAuthReq) (resp *types.GetTokenByAuthRes, err error) {
	// todo: add your logic here and delete this line

	return
}
