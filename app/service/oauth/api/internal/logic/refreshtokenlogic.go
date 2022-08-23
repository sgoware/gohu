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
	tokenGranter := token.GetTokenGranter()
	accessToken, err := tokenGranter.Grant(l.ctx, token.GrantByRefreshToken, req.RefreshToken)
	if err != nil {
		logx.Errorf("get token by auth failed, err: %v", err)
		return &types.GetTokenByRefreshTokenRes{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		}, nil
	}

	resp = &types.GetTokenByRefreshTokenRes{
		Code: http.StatusOK,
		Msg:  "get token successfully",
	}
	err = mapping.Struct2Struct(accessToken, &resp.Data.AccessToken)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
