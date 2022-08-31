package logic

import (
	"context"
	"main/app/service/user/rpc/crud/crud"

	"main/app/service/user/api/internal/svc"
	"main/app/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCollectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCollectionLogic {
	return &DeleteCollectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCollectionLogic) DeleteCollection(req *types.DeleteCollectionReq) (resp *types.DeleteCollectionRes, err error) {
	res, _ := l.svcCtx.CrudRpcClient.DeleteCollection(l.ctx, &crud.DeleteCollectionReq{CollectionId: req.CollectionId})
	return &types.DeleteCollectionRes{
		Code: res.Code,
		Msg:  res.Msg,
		Ok:   res.Ok,
	}, nil
}
