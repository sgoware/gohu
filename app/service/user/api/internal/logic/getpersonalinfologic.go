package logic

import (
	"context"
	"main/app/common/log"
	"main/app/service/user/rpc/info/info"
	"main/app/utils/mapping"
	"net/http"

	"main/app/service/user/api/internal/svc"
	"main/app/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPersonalInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPersonalInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPersonalInfoLogic {
	return &GetPersonalInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPersonalInfoLogic) GetPersonalInfo(req *types.GetPersonalInfoReq) (resp *types.GetPersonalInfoRes, err error) {
	logger := log.GetSugaredLogger()

	res, _ := l.svcCtx.InfoRpcClient.GetPersonalInfo(l.ctx, &info.GetPersonalInfoReq{
		UserId: req.UserId,
	})
	data := types.GetPersonalInfoResData{}
	err = mapping.Struct2Struct(res.Data, &data)
	if err != nil {
		logger.Errorf("mapping struct failed, err: %v", err)
		return &types.GetPersonalInfoRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
			Data: types.GetPersonalInfoResData{},
		}, nil
	}
	return &types.GetPersonalInfoRes{
		Code: res.Code,
		Msg:  res.Msg,
		Ok:   res.Ok,
		Data: data,
	}, nil
}
