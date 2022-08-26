package model

import (
	"main/app/common/config"
	"main/app/common/log"
	"time"

	"github.com/spf13/cast"
)

type UserDetail struct {
	UserId string `json:"user_id,omitempty"` // 用户 id
	//UserPassword string   //用户密码
	Username    string   `json:"username,omitempty"`    // 用户名(唯一)
	NickName    string   `json:"nick_name,omitempty"`   // 用户昵称
	LastIp      string   `json:"last_ip,omitempty"`     // 上次登录地址
	Vip         int      `json:"vip,omitempty"`         // vip 等级
	Status      int      `json:"status,omitempty"`      // 在线状态
	UpdateTime  string   `json:"update_time,omitempty"` // 用户信息更新时间
	CreateTime  string   `json:"create_time,omitempty"` // 用户信息创建时间
	Authorities []string `json:"authorities,omitempty"` // 用户具备的权限
}

type ClientDetail struct {
	ClientId                    string   `json:"client_id,omitempty"`                      // client 的标识
	AccessTokenValiditySeconds  int64    `json:"access_token_validity_seconds,omitempty"`  // 访问令牌有效时间，秒
	RefreshTokenValiditySeconds int64    `json:"refresh_token_validity_seconds,omitempty"` // 刷新令牌有效时间，秒
	RegisteredRedirectUri       string   `json:"registered_redirect_uri,omitempty"`        // 重定向地址，授权码类型中使用
	AuthorizedGrantTypes        []string `json:"authorized_grant_types,omitempty"`         // 可以使用的授权类型
}

type ClientDetailWithSecret struct {
	ClientSecret string `json:"client_secret,omitempty"` // client密钥
	ClientDetail
}

type OAuth2Token struct {
	RefreshToken *OAuth2Token `json:"refresh_token,omitempty"` // 刷新令牌
	TokenType    string       `json:"token_type,omitempty"`    // 令牌类型
	TokenValue   string       `json:"token_value,omitempty"`   // 令牌
	ExpiresAt    int64        `json:"expires_at,omitempty"`    // 过期时间
}

type OAuth2Details struct {
	Client *ClientDetail `json:"client,omitempty"` // 请求token的客户端信息
	User   *UserDetail   `json:"user,omitempty"`   // 请求token的用户信息
	Issuer string        `json:"issuer,omitempty"` // token签发者
}

const (
	JwtToken     = "jwt"
	AccessToken  = "access_token"
	RefreshToken = "refresh_token"
)

var (
	clientDetailsWithSecret map[string]ClientDetailWithSecret
	clientAuths             map[string]string
)

func (oauth2Token *OAuth2Token) IsExpired() bool {
	return oauth2Token.ExpiresAt != 0 && time.Unix(oauth2Token.ExpiresAt, 0).Before(time.Now())
}

func InitClientDetails() {
	logger := log.GetSugaredLogger()
	clientDetailsWithSecret = make(map[string]ClientDetailWithSecret)
	clientAuths = make(map[string]string)
	configClient, err := config.GetConfigClient()
	if err != nil {
		logger.Errorf("get configClient failed, err: %v", err)
	}
	mapIface := configClient.GetClientDetails()
	for k, val := range mapIface {
		details := cast.ToStringMap(val)
		var clientDetailWithSecret ClientDetailWithSecret
		clientDetailWithSecret.ClientId = k
		for detailName, detailValue := range details {
			switch detailName {
			case "secret":
				clientDetailWithSecret.ClientSecret = cast.ToString(detailValue)
				clientAuths[clientDetailWithSecret.ClientId] = clientDetailWithSecret.ClientSecret
			case "accesstokenvalidityseconds":
				clientDetailWithSecret.AccessTokenValiditySeconds = cast.ToInt64(detailValue)
			case "refreshtokenvalidityseconds":
				clientDetailWithSecret.RefreshTokenValiditySeconds = cast.ToInt64(detailValue)
			case "registeredredirecturi":
				clientDetailWithSecret.RegisteredRedirectUri = cast.ToString(detailValue)
			case "authorizedgranttypes":
				clientDetailWithSecret.AuthorizedGrantTypes = cast.ToStringSlice(detailValue)
			}
		}
		clientDetailsWithSecret[k] = clientDetailWithSecret
	}
	return
}

func GetClientDetails() map[string]ClientDetailWithSecret {
	return clientDetailsWithSecret
}

//func GetClientDetail() ClientDetailWithSecret {
//
//}
//
//func GetClientAuths() map[string]string {
//	return clientAuths
//}
