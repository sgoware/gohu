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

type GetObjInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetObjInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetObjInfoLogic {
	return &GetObjInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetObjInfoLogic) GetObjInfo(req *types.GetObjInfoReq) (resp *types.GetObjInfoRes, err error) {
	j := gjson.Parse(cast.ToString(l.ctx.Value("user_details")))
	userId := j.Get("user_id").Int()
	res, _ := l.svcCtx.InfoRpcClient.GetObjInfo(l.ctx, &info.GetObjInfoReq{
		UserId:  userId,
		ObjType: req.Obj_type,
	})
	return &types.GetObjInfoRes{
		Code: res.Code,
		Msg:  res.Msg,
		Ok:   res.Ok,
		Data: types.GetObjInfoResData{Ids: res.Data.Ids},
	}, nil
}
