package logic

import (
	"context"
	"fmt"
	"main/app/common/log"
	"main/app/service/oauth/model"
	"main/app/service/oauth/rpc/token/enhancer/internal/svc"
	"main/app/service/oauth/rpc/token/enhancer/pb"
	"main/app/service/oauth/rpc/token/store/tokenstore"
	"main/app/utils/mapping"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshAccessTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshAccessTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshAccessTokenLogic {
	return &RefreshAccessTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefreshAccessTokenLogic) RefreshAccessToken(in *pb.RefreshAccessTokenReq) (res *pb.RefreshAccessTokenRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	refreshToken, oauth2Details, err := l.svcCtx.Enhancer.ParseToken(in.RefreshToken)
	if err != nil {
		res = &pb.RefreshAccessTokenRes{
			Code: http.StatusOK,
			Msg:  fmt.Sprintf("parse oauth token failed, %v", err),
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
	logger.Debugf("refreshToken: \n%v, \noauth2Details: \n%v", refreshToken, oauth2Details)
	if time.Unix(refreshToken.ExpiresAt, 0).Before(time.Now()) {
		res = &pb.RefreshAccessTokenRes{
			Code: http.StatusOK,
			Msg:  model.ErrExpiredToken.Error(),
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	// TODO: 这里可以移除原有的访问令牌,也可以直接覆盖原有的访问令牌
	accessToken, err := l.svcCtx.Enhancer.GenerateToken(oauth2Details, model.AccessToken)
	if err != nil {
		res = &pb.RefreshAccessTokenRes{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Sprintf("refresh token failed, %v", err),
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
	res = &pb.RefreshAccessTokenRes{
		Code: http.StatusOK,
		Msg:  "refresh token successfully",
		Ok:   true,
		Data: &pb.RefreshAccessTokenRes_Data{AccessToken: &pb.OAuth2Token{}},
	}
	err = mapping.Struct2Struct(accessToken, res.Data.AccessToken)
	if err != nil {
		res = &pb.RefreshAccessTokenRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, err
	}

	// 存储刷新后的令牌到redis
	storeToken := &tokenstore.StoreTokenReq{
		UserId:      oauth2Details.User.UserId,
		AccessToken: &tokenstore.OAuth2Token{},
	}
	err = mapping.Struct2Struct(accessToken, storeToken.AccessToken)
	if err != nil {
		logger.Errorf("mapping struct failed, err: %v", err)
		res = &pb.RefreshAccessTokenRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
	storeTokenRes, _ := l.svcCtx.TokenStoreRpcClient.StoreToken(l.ctx, storeToken)
	if !storeTokenRes.Ok {
		res = &pb.RefreshAccessTokenRes{
			Code: res.Code,
			Msg:  fmt.Sprintf("store token failed, %v", storeTokenRes.Msg),
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	logger.Debugf("send message: %v", res.String())
	return res, nil
}
