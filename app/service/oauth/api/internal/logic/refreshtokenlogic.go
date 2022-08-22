package logic

import (
	"context"

	"gohu/app/service/oauth/api/internal/svc"
	"gohu/app/service/oauth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.GetTokenByRefreshTokenReq) (resp *types.GetTokenByRefreshTokenRes, err error) {
	// todo: add your logic here and delete this line

	return
}
