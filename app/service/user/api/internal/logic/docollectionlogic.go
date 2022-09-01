package logic

import (
	"context"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"main/app/service/user/rpc/crud/crud"

	"main/app/service/user/api/internal/svc"
	"main/app/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DoCollectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDoCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DoCollectionLogic {
	return &DoCollectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DoCollectionLogic) DoCollection(req *types.DoCollectionReq) (resp *types.DoCollectionRes, err error) {
	j := gjson.Parse(cast.ToString(l.ctx.Value("user_details")))
	userId := j.Get("user_id").Int()
	res, _ := l.svcCtx.CrudRpcClient.DoCollection(l.ctx, &crud.DoCollectionReq{
		UserId:      userId,
		CollectType: req.CollectionType,
		ObjType:     req.ObjType,
		ObjId:       req.ObjId,
	})
	return &types.DoCollectionRes{
		Code: res.Code,
		Msg:  res.Msg,
		Ok:   res.Ok,
	}, nil
}
