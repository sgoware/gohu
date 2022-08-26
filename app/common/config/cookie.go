package config

import (
	"net/http"
)

type CookieConfig struct {
	Secret string
	http.Cookie
}

func NewCookieConfig() (*CookieConfig, error) {
	if agolloClient == nil {
		return nil, errEmptyConfigClient
	}
	v, err := agolloClient.GetViper("oauth.yaml")
	if err != nil {
		return nil, errGetViper
	}
	return &CookieConfig{
		Secret: v.GetString("Cookie.Secret"),
		Cookie: http.Cookie{
			Domain:   v.GetString("Cookie.Domain"),
			MaxAge:   v.GetInt("Cookie.MaxAge"),
			Secure:   v.GetBool("Cookie.Secure"),
			HttpOnly: v.GetBool("Cookie.HttpOnly"),
		},
	}, nil
}
