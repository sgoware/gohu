package logic

import (
	"context"
	"gohu/app/common/log"
	"gohu/app/service/oauth/model"
	"gohu/app/service/oauth/rpc/token/enhancer/internal/svc"
	"gohu/app/service/oauth/rpc/token/enhancer/pb"
	"gohu/app/service/oauth/rpc/token/store/tokenstore"
	"gohu/app/utils/mapping"
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

	oauth2Details := &model.OAuth2Details{}
	err = mapping.Struct2Struct(in.Oauth2Details, oauth2Details)
	if err != nil {
		return nil, err
	}
	existTokenRes, _ := l.svcCtx.TokenStoreRpcClient.GetToken(l.ctx, &tokenstore.GetTokenReq{
		UserId: in.Oauth2Details.User.UserId,
	})
	if existTokenRes.Ok {
		if !time.Unix(existTokenRes.Data.OauthToken.ExpiresAt, 0).Before(time.Now()) {
			res = &pb.CreateAccessTokenRes{
				Ok:   true,
				Msg:  "create token successfully",
				Data: &pb.CreateAccessTokenRes_Data{AccessToken: &pb.OAuth2Token{}},
			}
			err = mapping.Struct2Struct(existTokenRes.Data.OauthToken, res.Data.AccessToken)
			if err != nil {
				return nil, err
			}
			return res, nil
		}
		// 访问令牌失效的情况,可以移除,也可以在后面创建新令牌的时候直接覆盖令牌
		// TODO:
	}
	accessToken, err := l.svcCtx.Enhancer.GenerateToken(oauth2Details, model.AccessToken)
	if err != nil {
		logger.Errorf("generate access_token failed, err: %v", err)
		return &pb.CreateAccessTokenRes{
			Ok:  false,
			Msg: "generate access_token failed",
		}, nil
	}
	res = &pb.CreateAccessTokenRes{
		Ok:   true,
		Msg:  "create access_token successfully",
		Data: &pb.CreateAccessTokenRes_Data{AccessToken: &pb.OAuth2Token{}},
	}
	err = mapping.Struct2Struct(accessToken, res.Data.AccessToken)
	if err != nil {
		return nil, err
	}
	storeToken := &tokenstore.StoreTokenReq{
		UserId:      oauth2Details.User.UserId,
		AccessToken: &tokenstore.OAuth2Token{},
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
