package logic

import (
	"context"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"main/app/service/user/rpc/info/info"

	"main/app/service/user/api/internal/svc"
	"main/app/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCollectionInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCollectionInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCollectionInfoLogic {
	return &GetCollectionInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCollectionInfoLogic) GetCollectionInfo(req *types.GetCollectionInfoReq) (resp *types.GetCollectionInfoRes, err error) {
	j := gjson.Parse(cast.ToString(l.ctx.Value("user_details")))
	userId := j.Get("user_id").Int()
	res, _ := l.svcCtx.InfoRpcClient.GetCollectionInfo(l.ctx, &info.GetCollectionInfoReq{
		UserId: userId,
	})
	return &types.GetCollectionInfoRes{
		Code: res.Code,
		Msg:  res.Msg,
		Ok:   res.Ok,
		Data: types.GetCollectionInfoResData{
			ObjType: res.Data.ObjType,
			ObjId:   res.Data.ObjId,
		},
	}, nil
}
