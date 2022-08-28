package model

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/spf13/cast"
	apollo "main/app/common/config"
	"main/app/service/user/dao/query"
	"strings"
)

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("default", "aertcerac", "744637972",ture).
func parseBasicAuth(auth string) (clientId, clientSecret, userId string, ok bool) {
	// TODO: 待测试 解码与字符串分割
	const prefix = "Basic "
	if len(auth) < len(prefix) || (auth[:len(prefix)] != prefix) {
		return "", "", "", false
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return "", "", "", false
	}
	cs := string(c)
	seq := strings.Split(cs, ":")
	if len(seq) == 3 {
		ok = true
	}
	clientId = seq[0]
	clientSecret = seq[1]
	userId = seq[2]
	if !ok {
		return "", "", "", false
	}
	return clientId, clientSecret, userId, true
}

type TokenGranter interface {
	Grant(ctx context.Context, grantType string, auth string) (*OAuth2Token, error)
}

type ComposeTokenGranter struct {
	TokenGrantDict map[string]TokenGranter
}

func NewComposeTokenGranter(tokenGrantDict map[string]TokenGranter) TokenGranter {
	return &ComposeTokenGranter{
		TokenGrantDict: tokenGrantDict,
	}
}

func (tokenGranter *ComposeTokenGranter) Grant(ctx context.Context, grantType string, auth string) (*OAuth2Token, error) {

	dispatchGranter := tokenGranter.TokenGrantDict[grantType]

	if dispatchGranter == nil {
		return nil, ErrNotSupportGrantType
	}

	return dispatchGranter.Grant(ctx, grantType, auth)
}

type AuthorizationTokenGranter struct {
	SupportGrantType string
	ClientDetails    map[string]ClientDetailWithSecret
	TokenService     TokenService
	UserModel        *query.Query
}

func NewAuthorizationTokenGranter(grantType string, clientDetails map[string]ClientDetailWithSecret, tokenService TokenService) (TokenGranter, error) {
	if grantType == "" || clientDetails == nil || tokenService == nil {
		return nil, errors.New("param cannot be null")
	}
	db, err := apollo.GetMysqlDB("user.yaml")
	if err != nil {
		return nil, errors.New("initialize mysql failed")
	}
	return &AuthorizationTokenGranter{
		SupportGrantType: grantType,
		ClientDetails:    clientDetails,
		TokenService:     tokenService,
		UserModel:        query.Use(db),
	}, nil
}

func (tokenGranter *AuthorizationTokenGranter) Grant(ctx context.Context,
	grantType string, auth string) (*OAuth2Token, error) {

	if grantType != tokenGranter.SupportGrantType {
		return nil, ErrNotSupportGrantType
	}
	// 匹配body中的authorization
	clientId, clientSecret, userId, ok := parseBasicAuth(auth)
	if !ok || clientSecret != tokenGranter.ClientDetails[clientId].ClientSecret {
		return nil, ErrInvalidAuthorizationRequest
	}

	userSubjectModel := tokenGranter.UserModel.UserSubject
	userDetail, err := userSubjectModel.WithContext(context.Background()).
		Where(userSubjectModel.ID.Eq(cast.ToInt64(userId))).First()
	if err != nil {
		return nil, ErrUserDetailNotFound
	}

	// 根据用户信息和客户端信息生成访问令牌
	return tokenGranter.TokenService.CreateAccessToken(ctx, &OAuth2Details{
		Client: &ClientDetail{
			ClientId:                    clientId,
			AccessTokenValiditySeconds:  tokenGranter.ClientDetails[clientId].AccessTokenValiditySeconds,
			RefreshTokenValiditySeconds: tokenGranter.ClientDetails[clientId].RefreshTokenValiditySeconds,
			RegisteredRedirectUri:       tokenGranter.ClientDetails[clientId].RegisteredRedirectUri,
			AuthorizedGrantTypes:        tokenGranter.ClientDetails[clientId].AuthorizedGrantTypes,
		},

		User: &UserDetail{
			UserId:      userDetail.ID,
			Username:    userDetail.Username,
			NickName:    userDetail.Nickname,
			LastIp:      userDetail.LastIP,
			Vip:         userDetail.Vip,
			State:       userDetail.State,
			UpdateTime:  userDetail.UpdateTime.Format("2006-01-02 15:04:05 MST"),
			CreateTime:  userDetail.CreateTime.Format("2006-01-02 15:04:05 MST"),
			Authorities: nil, // TODO: 待加入权限字段
		},
	})

}

type RefreshTokenGranter struct {
	SupportGrantType string
	TokenService     TokenService
}

func NewRefreshGranter(grantType string, tokenService TokenService) (TokenGranter, error) {
	if grantType == "" || tokenService == nil {
		return nil, errors.New("param cannot be null")
	}
	return &RefreshTokenGranter{
		SupportGrantType: grantType,
		TokenService:     tokenService,
	}, nil
}

func (tokenGranter *RefreshTokenGranter) Grant(ctx context.Context, grantType, token string) (*OAuth2Token, error) {
	if grantType != tokenGranter.SupportGrantType {
		return nil, ErrNotSupportGrantType
	}
	if token == "" {
		return nil, ErrInvalidToken
	}
	accessToken, err := tokenGranter.TokenService.RefreshAccessToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}
