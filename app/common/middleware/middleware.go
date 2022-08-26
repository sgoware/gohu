package middleware

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/imroc/req/v3"
	"github.com/spf13/cast"
	"github.com/thedevsaddam/gojsonq/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"main/app/common/config"
	"main/app/common/model/response"
	"main/app/service/oauth/model"
	"main/app/utils/cookie"
	"main/app/utils/jwt"
	"net/http"
)

type AuthMiddleware struct {
	Domain string
	*config.CookieConfig
	Rdb *redis.Client
}

func NewAuthMiddleware(domain string, cookieConfig *config.CookieConfig, rdb *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{Domain: domain, CookieConfig: cookieConfig, Rdb: rdb}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		var accessToken string
		var refreshToken string
		cookieWriter := cookie.NewCookieWriter(m.CookieConfig.Secret, cookie.Option{
			Config:  m.Cookie,
			Writer:  w,
			Request: r,
		})
		ok := cookieWriter.Get("x-token", &accessToken)
		if accessToken == "" || !ok {
			response.ResultWithData(w, http.StatusForbidden, "not logged in or illegal access", map[string]interface{}{"reload": true})
			return
		}

		if jwt.IsBlacklist(m.Rdb, accessToken) {
			response.ResultWithData(w, http.StatusForbidden, "login from a different location or accessToken exceeded", map[string]interface{}{"reload": true})
			return
		}
		// TODO: 使用oauth2服务器认证,使用认证令牌认证,如果过期则使用刷新令牌
		res, err := req.NewRequest().SetFormData(map[string]string{"oauth2_token": accessToken, "token_type": model.AccessToken}).
			Post("https://" + m.Domain + "/api/oauth/token/check")
		if err != nil {
			logx.Errorf("%v", err)
			return
		}
		if res.StatusCode != http.StatusOK {
			logx.Errorf("%v", res)
			return
		}
		ok = cast.ToBool(gojsonq.New().FromString(res.String()).Find("ok"))
		if !ok {
			//不ok则认证失败，包括刷新令牌
			//认证的时候若认证令牌过期，则刷新令牌
			msg := cast.ToString(gojsonq.New().FromString(res.String()).Find("msg"))
			// 认证令牌过期,用刷新令牌刷新
			if msg == "accessToken is expired" {
				ok = cookieWriter.Get("refresh-token", &refreshToken)
				if accessToken == "" || !ok {
					response.ResultWithData(w, http.StatusForbidden, "illegal access", map[string]interface{}{"reload": true})
					return
				}
				res, err = req.NewRequest().SetPathParam("refresh-token", refreshToken).
					Post("https://" + m.Domain + "/api/oauth/token/refresh")
				if err != nil {
					logx.Errorf("%v", err)
					return
				}
				if res.StatusCode != http.StatusOK {
					return
				}
				accessToken = cast.ToString(
					gojsonq.New().FromString(res.String()).
						Find("data.access_token.token_value"))
				refreshToken = cast.ToString(
					gojsonq.New().FromString(res.String()).
						Find("data.access_token.refresh_token.token_value"))
				cookieWriter.Set("x-token", accessToken)
				cookieWriter.Set("refresh-token", refreshToken)
			}
		}
		userDetailsRes, err := req.NewRequest().SetFormData(map[string]string{"access_token": accessToken}).
			Post("https://" + m.Domain + "/api/oauth/token/get/user")
		if err != nil {
			logx.Errorf("%v", err)
			return
		}
		if userDetailsRes.StatusCode != http.StatusOK {
			logx.Errorf("%v", res)
			return
		}
		j := gojsonq.New().FromString(userDetailsRes.String())
		// TODO: 待添加信息
		userDetails := &model.UserDetail{
			UserId:      cast.ToString(j.Find("data.user_details.user_id")),
			Username:    "",
			LastIp:      "",
			Status:      0,
			UpdateTime:  "",
			CreateTime:  "",
			Authorities: nil,
		}
		r = r.WithContext(context.WithValue(r.Context(), "user_id", userDetails.UserId))
		next(w, r)
	}
}
