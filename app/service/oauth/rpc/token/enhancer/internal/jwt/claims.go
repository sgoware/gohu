package jwt

import (
	"gohu/app/service/oauth/model"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	jwt.RegisteredClaims                   // token注册信息
	RefreshToken         model.OAuth2Token // 刷新令牌
	BaseClaims
}

type BaseClaims struct {
	model.UserDetail
	model.ClientDetail
}
