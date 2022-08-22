package logic

import (
	"context"
	"gohu/app/service/oauth/api/internal/svc"
	"gohu/app/service/oauth/api/internal/token"
	"gohu/app/service/oauth/api/internal/types"
	"gohu/app/utils/mapping"
	"net/http"

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

func (l *ReadTokenLogic) ReadToken(req *types.ReadTokenReq) (*types.ReadTokenRes, error) {
	tokenService := token.GetTokenService()
	oauthToken, err := tokenService.ReadAccessToken(l.ctx, req.OAuth2Token)
	if err != nil {
		logx.Errorf("parse token failed, err: %v", err)
		return &types.ReadTokenRes{
			Code: http.StatusBadRequest,
			Msg:  "invalid token",
		}, nil
	}

	resp := &types.ReadTokenRes{
		Code: http.StatusOK,
		Msg:  "check token successfully",
		Data: types.ReadTokenResData{AccessToken: &types.OAuth2Token{}},
	}
	err = mapping.Struct2Struct(oauthToken, resp.Data.AccessToken)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
