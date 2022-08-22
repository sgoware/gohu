package model

import (
	"context"
	"encoding/base64"
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
}

func NewAuthorizationTokenGranter(grantType string, clientDetails map[string]ClientDetailWithSecret, tokenService TokenService) TokenGranter {
	return &AuthorizationTokenGranter{
		SupportGrantType: grantType,
		ClientDetails:    clientDetails,
		TokenService:     tokenService,
	}
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

	// 根据用户信息和客户端信息生成访问令牌
	return tokenGranter.TokenService.CreateAccessToken(ctx, &OAuth2Details{
		Client: &ClientDetail{
			ClientId:                    clientId,
			AccessTokenValiditySeconds:  tokenGranter.ClientDetails[clientId].AccessTokenValiditySeconds,
			RefreshTokenValiditySeconds: tokenGranter.ClientDetails[clientId].RefreshTokenValiditySeconds,
			RegisteredRedirectUri:       tokenGranter.ClientDetails[clientId].RegisteredRedirectUri,
			AuthorizedGrantTypes:        tokenGranter.ClientDetails[clientId].AuthorizedGrantTypes,
		},
		User: &UserDetail{UserId: userId},
	})

}

type RefreshTokenGranter struct {
	SupportGrantType string
	TokenService     TokenService
}

func NewRefreshGranter(grantType string, tokenService TokenService) TokenGranter {
	return &RefreshTokenGranter{
		SupportGrantType: grantType,
		TokenService:     tokenService,
	}
}

func (tokenGranter *RefreshTokenGranter) Grant(ctx context.Context, grantType, token string) (*OAuth2Token, error) {
	if grantType != tokenGranter.SupportGrantType {
		return nil, ErrNotSupportGrantType
	}
	if token == "" {
		return nil, ErrInvalidTokenRequest
	}
	accessToken, err := tokenGranter.TokenService.RefreshAccessToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}
