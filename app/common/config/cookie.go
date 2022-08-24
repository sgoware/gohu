package config

import "net/http"

type CookieConfig struct {
	Secret string
	http.Cookie
}

func (c *Agollo) NewCookieConfig() (*CookieConfig, error) {
	v, err := c.GetViper("oauth.yaml")
	if err != nil {
		return nil, err
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
