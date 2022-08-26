package config

import (
	"errors"
	"net/http"
)

type CookieConfig struct {
	Secret string
	http.Cookie
}

func NewCookieConfig() (*CookieConfig, error) {
	if agolloClient == nil {
		return nil, errors.New(emptyConfigClientErr)
	}
	v, err := agolloClient.GetViper("oauth.yaml")
	if err != nil {
		return nil, errors.New(getViperErr)
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
