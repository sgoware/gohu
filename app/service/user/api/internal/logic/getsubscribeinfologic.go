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

type GetSubscribeInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSubscribeInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSubscribeInfoLogic {
	return &GetSubscribeInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSubscribeInfoLogic) GetSubscribeInfo(req *types.GetSubscribeInfoReq) (resp *types.GetSubscribeInfoRes, err error) {
	j := gjson.Parse(cast.ToString(l.ctx.Value("user_details")))
	userId := j.Get("user_id").Int()
	res, _ := l.svcCtx.InfoRpcClient.GetSubscribeInfo(l.ctx, &info.GetSubscribeInfoReq{
		UserId:  userId,
		ObjType: req.ObjType,
	})
	return &types.GetSubscribeInfoRes{
		Code: res.Code,
		Msg:  res.Msg,
		Ok:   res.Ok,
		Data: types.GetSubscribeInfoResData{
			Ids: res.Data.Ids,
		},
	}, nil
}
