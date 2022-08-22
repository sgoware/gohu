package model

import (
	"context"
	"gohu/app/service/oauth/rpc/token/enhancer/tokenenhancer"
	"gohu/app/service/oauth/rpc/token/store/tokenstore"
	"gohu/app/utils/mapping"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type TokenService interface {
	// CreateAccessToken 根据用户信息和客户端信息生成访问令牌
	CreateAccessToken(ctx context.Context, oauth2Details *OAuth2Details) (*OAuth2Token, error)
	// RefreshAccessToken 根据刷新令牌获取访问令牌
	RefreshAccessToken(ctx context.Context, refreshToken string) (*OAuth2Token, error)
	// ReadAccessToken 根据访问令牌值获取访问令牌结构体
	ReadAccessToken(ctx context.Context, accessToken string) (*OAuth2Token, error)
	// GetUserDetails 获取令牌对应的用户信息
	GetUserDetails(ctx context.Context, accessToken string) (*UserDetail, error)

	//// GetOAuth2DetailsByAccessToken 根据访问令牌获取对应的用户信息和客户端信息
	//GetOAuth2DetailsByAccessToken(tokenValue string) (*OAuth2Details, error)
	//// GetAccessToken 根据用户信息和客户端信息获取已生成访问令牌
	//GetAccessToken(details *OAuth2Details) (*OAuth2Token, error)
}

type DefaultTokenService struct {
}

type RpcTokenService struct {
	TokenEnhancerClient tokenenhancer.TokenEnhancer
	TokenStoreClient    tokenstore.TokenStore
}

func NewRpcTokenService(enhancerConf, storeConf zrpc.RpcClientConf) TokenService {
	return &RpcTokenService{
		TokenEnhancerClient: tokenenhancer.NewTokenEnhancer(zrpc.MustNewClient(enhancerConf)),
		TokenStoreClient:    tokenstore.NewTokenStore(zrpc.MustNewClient(storeConf)),
	}
}

func (tokenService *RpcTokenService) CreateAccessToken(ctx context.Context, oauth2Details *OAuth2Details) (*OAuth2Token, error) {
	res, err := tokenService.TokenEnhancerClient.CreateAccessToken(ctx, &tokenenhancer.CreateAccessTokenReq{
		Oauth2Details: &tokenenhancer.OAuth2Details{
			Client: &tokenenhancer.ClientDetails{
				ClientId:                    oauth2Details.Client.ClientId,
				AccessTokenValiditySeconds:  oauth2Details.Client.AccessTokenValiditySeconds,
				RefreshTokenValiditySeconds: oauth2Details.Client.RefreshTokenValiditySeconds,
				RegisteredRedirectUri:       oauth2Details.Client.RegisteredRedirectUri,
				AuthorizedGrantTypes:        oauth2Details.Client.AuthorizedGrantTypes,
			},
			User: &tokenenhancer.UserDetails{
				UserId: oauth2Details.User.UserId,
			},
		},
	})
	if err != nil {
		logx.Errorf("create access_token failed, err: %v", err)
		return nil, err
	}
	accessToken := &OAuth2Token{}
	err = mapping.Struct2Struct(res.Data.AccessToken, accessToken)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (tokenService *RpcTokenService) RefreshAccessToken(ctx context.Context, refreshTokenValue string) (*OAuth2Token, error) {
	res, err := tokenService.TokenEnhancerClient.RefreshAccessToken(ctx,
		&tokenenhancer.RefreshAccessTokenReq{RefreshToken: refreshTokenValue})
	if err != nil {
		logx.Errorf("refresh access_token failed, err: %v", err)
		return nil, err
	}
	accessToken := &OAuth2Token{}
	err = mapping.Struct2Struct(res.Data.AccessToken, accessToken)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (tokenService *RpcTokenService) ReadAccessToken(ctx context.Context, accessTokenValue string) (*OAuth2Token, error) {
	res, err := tokenService.TokenEnhancerClient.ReadOauthToken(ctx, &tokenenhancer.ReadTokenReq{OauthToken: accessTokenValue})
	if err != nil {
		logx.Errorf("read access_token failed, err: %v", err)
		return nil, err
	}
	accessToken := &OAuth2Token{}
	err = mapping.Struct2Struct(res.Data.AccessToken, accessToken)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (tokenService *RpcTokenService) GetUserDetails(ctx context.Context, accessTokenValue string) (*UserDetail, error) {
	res, err := tokenService.TokenEnhancerClient.GetUserDetails(ctx, &tokenenhancer.GetUserDetailsReq{AccessToken: accessTokenValue})
	if err != nil {
		logx.Errorf("get user details failed, err: %v", err)
		return nil, err
	}
	userDetails := &UserDetail{}
	err = mapping.Struct2Struct(res.Data.UserDetails, userDetails)
	if err != nil {
		return nil, err
	}
	return userDetails, nil
}
