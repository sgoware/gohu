package logic

import (
	"context"

	"gohu/app/service/oauth/api/internal/svc"
	"gohu/app/service/oauth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadTokenLogic {
	return &ReadTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadTokenLogic) ReadToken(req *types.ReadTokenReq) (resp *types.ReadTokenRes, err error) {
	// todo: add your logic here and delete this line

	return
}
