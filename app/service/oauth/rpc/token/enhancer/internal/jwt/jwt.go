package jwt

import (
	"errors"
	"gohu/app/service/oauth/model"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cast"
)

type JWT struct {
	SigningKey []byte
	Issuer     string
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

func NewJWT(signingKey, issuer string) *JWT {
	return &JWT{
		SigningKey: []byte(signingKey),
		Issuer:     issuer,
	}
}

func (j *JWT) GenerateToken(oauth2Details *model.OAuth2Details, tokenType string) (*model.OAuth2Token, error) {
	claims := j.createClaims(oauth2Details, tokenType)
	if tokenType == model.AccessToken {
		refreshToken, err := j.GenerateToken(oauth2Details, model.RefreshToken)
		if err != nil {
			return nil, err
		}
		claims.RefreshToken = *refreshToken
	}
	tokenOperator := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenValue, err := tokenOperator.SignedString(j.SigningKey)
	return &model.OAuth2Token{
		RefreshToken: &claims.RefreshToken,
		TokenType:    tokenType,
		TokenValue:   tokenValue,
		ExpiresAt:    claims.ExpiresAt.Time.Unix(),
	}, err
}

// ParseToken 根据访问令牌返回刷新令牌,用户信息,客户端信息
func (j *JWT) ParseToken(tokenValue string) (*model.OAuth2Token, *model.OAuth2Details, error) {
	tokenOperator, err := jwt.ParseWithClaims(tokenValue, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if v, ok := err.(*jwt.ValidationError); ok {
			if v.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, nil, TokenMalformed
			} else if v.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, nil, TokenExpired
			} else if v.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, nil, TokenNotValidYet
			} else {
				return nil, nil, TokenInvalid
			}
		}
	}
	if tokenOperator != nil {
		if claims, ok := tokenOperator.Claims.(*CustomClaims); ok && tokenOperator.Valid {
			if !reflect.DeepEqual(claims.RefreshToken, model.OAuth2Token{}) {
				return &model.OAuth2Token{
						RefreshToken: &claims.RefreshToken,
						TokenType:    model.AccessToken,
						TokenValue:   tokenValue,
						ExpiresAt:    claims.ExpiresAt.Time.Unix(),
					}, &model.OAuth2Details{
						Client: &claims.ClientDetail,
						User:   &claims.UserDetail,
						Issuer: claims.Issuer,
					}, nil
			} else {
				return &model.OAuth2Token{
						RefreshToken: nil,
						TokenType:    model.RefreshToken,
						TokenValue:   tokenValue,
						ExpiresAt:    claims.ExpiresAt.Time.Unix(),
					}, &model.OAuth2Details{
						Client: &claims.ClientDetail,
						User:   &claims.UserDetail,
						Issuer: claims.Issuer,
					}, nil

			}

		}
		return nil, nil, TokenInvalid
	} else {
		return nil, nil, TokenInvalid
	}
}

func (j *JWT) GetUserDetails(tokenValue string) (*model.UserDetail, error) {
	tokenOperator, err := jwt.ParseWithClaims(tokenValue, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if v, ok := err.(*jwt.ValidationError); ok {
			if v.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if v.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if v.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if tokenOperator != nil {
		if claims, ok := tokenOperator.Claims.(*CustomClaims); ok && tokenOperator.Valid {
			if !reflect.DeepEqual(claims.RefreshToken, model.OAuth2Token{}) {
				return &model.UserDetail{
					UserId:      claims.UserId,
					Username:    claims.Username,
					LastIp:      claims.LastIp,
					Status:      claims.Status,
					UpdateTime:  claims.UpdateTime,
					CreateTime:  claims.CreateTime,
					Authorities: claims.Authorities,
				}, nil
			} else {
				return nil, TokenInvalid

			}

		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}

func (j *JWT) createClaims(oauth2Details *model.OAuth2Details, tokenType string) CustomClaims {
	var validitySecond time.Duration
	if tokenType == model.AccessToken {
		validitySecond, _ = time.ParseDuration(cast.ToString(oauth2Details.Client.AccessTokenValiditySeconds) + "s")
	} else {
		validitySecond, _ = time.ParseDuration(cast.ToString(oauth2Details.Client.RefreshTokenValiditySeconds) + "s")
	}
	expiresAt := time.Now().Add(validitySecond)
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.Issuer,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(time.Now().Truncate(time.Second)),
		},
		BaseClaims: BaseClaims{
			UserDetail:   *oauth2Details.User,
			ClientDetail: *oauth2Details.Client,
		},
	}
	return claims
}
