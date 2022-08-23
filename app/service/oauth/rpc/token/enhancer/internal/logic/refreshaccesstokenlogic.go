package logic

import (
	"context"
	"main/app/common/log"
	"main/app/service/oauth/model"
	"main/app/service/oauth/rpc/token/enhancer/internal/svc"
	"main/app/service/oauth/rpc/token/enhancer/pb"
	"main/app/service/oauth/rpc/token/store/tokenstore"
	"main/app/utils/mapping"
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

	refreshToken, oauth2Details, err := l.svcCtx.Enhancer.ParseToken(in.RefreshToken)
	if err != nil {
		return &pb.RefreshAccessTokenRes{
			Ok:  false,
			Msg: model.ErrInvalidTokenRequest.Error(),
		}, nil
	}
	if time.Unix(refreshToken.ExpiresAt, 0).Before(time.Now()) {
		return &pb.RefreshAccessTokenRes{
			Ok:  false,
			Msg: model.ErrExpiredToken.Error(),
		}, nil
	}
	// TODO: 这里可以移除原有的访问令牌,也可以直接覆盖原有的访问令牌
	accessToken, err := l.svcCtx.Enhancer.GenerateToken(oauth2Details, model.AccessToken)
	if err != nil {
		logger.Errorf("refresh token failed, err: %v", err)
		return &pb.RefreshAccessTokenRes{
			Ok:  false,
			Msg: "refresh token failed",
		}, nil
	}
	res = &pb.RefreshAccessTokenRes{
		Ok:  true,
		Msg: "refresh token successfully",
	}
	err = mapping.Struct2Struct(accessToken, res.Data.AccessToken)
	if err != nil {
		return nil, err
	}
	storeToken := &tokenstore.StoreTokenReq{
		UserId: oauth2Details.User.UserId,
	}
	err = mapping.Struct2Struct(accessToken, storeToken.AccessToken)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.TokenStoreRpcClient.StoreToken(l.ctx, storeToken)
	if err != nil {
		return nil, err
	}
	return res, nil
}
