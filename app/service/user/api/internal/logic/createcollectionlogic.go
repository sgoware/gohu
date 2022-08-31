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

type CreateCollectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCollectionLogic {
	return &CreateCollectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCollectionLogic) CreateCollection(req *types.CreateCollectionReq) (resp *types.CreateCollectionRes, err error) {
	j := gjson.Parse(cast.ToString(l.ctx.Value("user_details")))
	userId := j.Get("user_id").Int()
	res, _ := l.svcCtx.CrudRpcClient.CreateCollection(l.ctx, &crud.CreateCollectionReq{
		UserId:  userId,
		ObjType: req.ObjType,
		ObjId:   req.ObjId,
	})
	return &types.CreateCollectionRes{
		Code: res.Code,
		Msg:  res.Msg,
		Ok:   res.Ok,
	}, nil
}
