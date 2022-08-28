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

type CreateAccessTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAccessTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAccessTokenLogic {
	return &CreateAccessTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAccessTokenLogic) CreateAccessToken(in *pb.CreateAccessTokenReq) (res *pb.CreateAccessTokenRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	oauth2Details := &model.OAuth2Details{}
	err = mapping.Struct2Struct(in.Oauth2Details, oauth2Details)
	if err != nil {
		res = &pb.CreateAccessTokenRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
	existTokenRes, _ := l.svcCtx.TokenStoreRpcClient.GetToken(l.ctx, &tokenstore.GetTokenReq{
		UserId: in.Oauth2Details.User.UserId,
	})
	if existTokenRes.Ok {
		// 获取到存在的令牌
		if !time.Unix(existTokenRes.Data.OauthToken.ExpiresAt, 0).Before(time.Now()) {
			res = &pb.CreateAccessTokenRes{
				Ok:   true,
				Msg:  "create token successfully",
				Data: &pb.CreateAccessTokenRes_Data{AccessToken: &pb.OAuth2Token{}},
			}
			err = mapping.Struct2Struct(existTokenRes.Data.OauthToken, res.Data.AccessToken)
			if err != nil {
				logger.Errorf("mapping struct failed, err: %v", err)
				res = &pb.CreateAccessTokenRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %v", res.String())
				return res, nil
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
	}
	// 生成新的令牌
	accessToken, err := l.svcCtx.Enhancer.GenerateToken(oauth2Details, model.AccessToken)
	if err != nil {
		res = &pb.CreateAccessTokenRes{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Sprintf("generate access_token failed, %v", err),
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
	logger.Debugf("generated token: %v", accessToken)
	res = &pb.CreateAccessTokenRes{
		Ok:   true,
		Msg:  "create access_token successfully",
		Data: &pb.CreateAccessTokenRes_Data{AccessToken: &pb.OAuth2Token{}},
	}
	err = mapping.Struct2Struct(accessToken, res.Data.AccessToken)
	if err != nil {
		logger.Errorf("mapping struct failed, err: %v", err)
		res = &pb.CreateAccessTokenRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	// 存储令牌到redis
	storeToken := &tokenstore.StoreTokenReq{
		UserId:      oauth2Details.User.UserId,
		AccessToken: &tokenstore.OAuth2Token{},
	}
	err = mapping.Struct2Struct(accessToken, storeToken.AccessToken)
	if err != nil {
		logger.Errorf("mapping struct failed, err: %v", err)
		res = &pb.CreateAccessTokenRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
	storeTokenRes, _ := l.svcCtx.TokenStoreRpcClient.StoreToken(l.ctx, storeToken)
	if !storeTokenRes.Ok {
		res = &pb.CreateAccessTokenRes{
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
