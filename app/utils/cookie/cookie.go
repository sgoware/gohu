package cookie

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type (
	Cookie struct {
		Secret string
		Opt    Option
	}
	Option struct {
		Config  http.Cookie
		Writer  http.ResponseWriter
		Request *http.Request
	}
)

func NewCookieWriter(secret string, opt ...Option) *Cookie {
	if len(opt) == 0 {
		return &Cookie{Secret: secret}
	} else {
		return &Cookie{
			Secret: secret,
			Opt:    opt[0],
		}
	}
}

// Set 写入数据的方法
func (c *Cookie) Set(key string, value interface{}) {
	bytes, _ := json.Marshal(value)
	setSecureCookie(c, key, string(bytes))
}

// Get 获取数据的方法
func (c *Cookie) Get(key string, obj interface{}) bool {
	tempData, ok := getSecureCookie(c, key)
	if !ok {
		return false
	}
	_ = json.Unmarshal([]byte(tempData), obj)
	return true
}

// Remove 删除数据的方法
func (c *Cookie) Remove(key string, value interface{}) {
	bytes, _ := json.Marshal(value)
	setSecureCookie(c, key, string(bytes))
}

func setSecureCookie(c *Cookie, name, value string) {
	vs := base64.URLEncoding.EncodeToString([]byte(value))
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := hmac.New(sha256.New, []byte(c.Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)

	sig := fmt.Sprintf("%02x", h.Sum(nil))
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")

	http.SetCookie(c.Opt.Writer, &http.Cookie{
		Name:     name,
		Value:    cookie,
		MaxAge:   c.Opt.Config.MaxAge,
		Path:     "/",
		Domain:   c.Opt.Config.Domain,
		SameSite: http.SameSite(1),
		Secure:   c.Opt.Config.Secure,
		HttpOnly: c.Opt.Config.HttpOnly,
	})
}

// GetSecureCookie Get secure cookie from request by a given key.
func getSecureCookie(c *Cookie, key string) (string, bool) {
	cookie, err := c.Opt.Request.Cookie(key)
	if err != nil {
		logx.Errorf("cookies not found, err: %v", err)
		return "", false
	}
	val, err := url.QueryUnescape(cookie.Value)
	if val == "" || err != nil {
		return "", false
	}

	parts := strings.SplitN(val, "|", 3)
	if len(parts) != 3 {
		return "", false
	}

	vs := parts[0]
	timestamp := parts[1]
	sig := parts[2]

	h := hmac.New(sha256.New, []byte(c.Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)

	if fmt.Sprintf("%02x", h.Sum(nil)) != sig {
		return "", false
	}
	res, _ := base64.URLEncoding.DecodeString(vs)
	return string(res), true
}
