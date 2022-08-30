package logic

import (
	"context"
	"fmt"
	"main/app/common/log"
	"main/app/service/oauth/api/internal/svc"
	"main/app/service/oauth/api/internal/token"
	"main/app/service/oauth/api/internal/types"
	"main/app/utils/mapping"
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
	logger := log.GetSugaredLogger()

	tokenService := token.GetTokenService()
	oauthToken, err := tokenService.ReadAccessToken(l.ctx, req.OAuth2Token)
	if err != nil {
		return &types.ReadTokenRes{
			Code: http.StatusBadRequest,
			Msg:  fmt.Sprintf("invalid token, %v", err),
			Ok:   false,
		}, nil
	}

	resp := &types.ReadTokenRes{
		Code: http.StatusOK,
		Msg:  "check token successfully",
		Ok:   true,
		Data: types.ReadTokenResData{AccessToken: &types.OAuth2Token{}},
	}
	err = mapping.Struct2Struct(oauthToken, resp.Data.AccessToken)
	if err != nil {
		logger.Errorf("mapping struct failed, err: %v", err)
		return &types.ReadTokenRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}, nil
	}
	return resp, nil
}
