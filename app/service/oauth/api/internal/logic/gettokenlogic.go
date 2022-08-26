package logic

import (
	"context"
	"main/app/service/oauth/api/internal/svc"
	"main/app/service/oauth/api/internal/token"
	"main/app/service/oauth/api/internal/types"
	"main/app/utils/mapping"
	"net/http"

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

func (l *GetTokenLogic) GetToken(req *types.GetTokenByAuthReq) (res *types.GetTokenByAuthRes, err error) {
	tokenGranter := token.GetTokenGranter()
	accessToken, err := tokenGranter.Grant(l.ctx, token.GrantByAuth, req.Authorization)
	if err != nil {
		logx.Errorf("get token by auth failed, err: %v", err)
		return &types.GetTokenByAuthRes{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		}, nil
	}

	res = &types.GetTokenByAuthRes{
		Code: http.StatusOK,
		Msg:  "get token successfully",
		Data: types.GetTokenByAuthResData{AccessToken: &types.OAuth2Token{}},
	}
	err = mapping.Struct2Struct(accessToken, res.Data.AccessToken)
	if err != nil {
		return nil, err
	}
	return res, nil
}
