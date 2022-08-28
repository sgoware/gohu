package logic

import (
	"context"
	"main/app/common/log"
	"main/app/service/oauth/api/internal/svc"
	"main/app/service/oauth/api/internal/token"
	"main/app/service/oauth/api/internal/types"
	"main/app/utils/mapping"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoReq) (*types.GetUserInfoRes, error) {
	logger := log.GetSugaredLogger()

	tokenService := token.GetTokenService()
	userDetails, err := tokenService.GetUserDetails(l.ctx, req.AccessToken)
	if err != nil {
		return &types.GetUserInfoRes{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		}, nil
	}
	resp := &types.GetUserInfoRes{
		Code: http.StatusOK,
		Msg:  "get user details successfully",
		Data: types.GetUserInfoResData{UserDetails: &types.UserDetails{}},
	}
	err = mapping.Struct2Struct(userDetails, resp.Data.UserDetails)
	if err != nil {
		logger.Errorf("mapping struct failed, err: %v", err)
		return &types.GetUserInfoRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}
	return resp, nil
}
