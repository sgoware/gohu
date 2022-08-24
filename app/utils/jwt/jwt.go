package jwt

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
)

type JWT struct {
	SigningKey []byte
}

type BlackList struct {
	Id  string
	Jwt string
}

type Config struct {
	SecretKey   string // 密钥
	ExpiresTime int64  // 过期时间,单位:秒
	BufferTime  int64  // 缓冲时间,缓冲时间内会获得新的token刷新令牌,此时一个用户会存在两个有效令牌,但是前端只留一个,另一个会丢失
	Issuer      string // 签发者
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

func NewJWT() *JWT {
	return &JWT{}
}

func (j *JWT) CreateClaims(config Config, baseClaims BaseClaims) CustomClaims {
	j.SigningKey = []byte(config.SecretKey)
	claims := CustomClaims{
		BufferTime: config.BufferTime,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Truncate(time.Second)), // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.ExpiresTime) * time.Second)),
			Issuer:    config.Issuer,
		},
		BaseClaims: baseClaims,
	}
	return claims
}

func (j *JWT) GenerateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//func (j *JWT) CreateTokenByOldToken(oldToken string, claims CustomClaims) (string, error) {
//	v, err, _ := g.ConcurrencyControl.Do("JWT_"+oldToken, func() (interface{}, error) {
//		return j.GenerateToken(claims)
//	})
//	return v.(string), err
//}

// ParseToken 解析JWT
func (j *JWT) ParseToken(secret, tokenString string) (*CustomClaims, error) {
	// 解析token
	j.SigningKey = []byte(secret)
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}

func IsBlacklist(gdb *redis.Client, jwt string) bool {
	val, err := gdb.Get(context.Background(), jwt).Result()
	if err != nil || val == "" {
		return false
	}
	return true
	// err := global.GVA_DB.Where("jwt = ?", jwt).First(&system.JwtBlacklist{}).Error
	// isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	// return !isNotFound
}

func (j *JWT) JsonInBlackList(rdb *redis.Client, jwtStr BlackList) error {
	err := rdb.Set(context.Background(), "jwt_"+jwtStr.Id, jwtStr.Jwt, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
