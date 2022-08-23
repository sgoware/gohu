// Code generated by goctl. DO NOT EDIT.
package types

type RegisterReq struct {
	Uid      string `form:"uid"`
	Nickname string `form:"nickname"`
	Password string `form:"password"`
}

type RegisterRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type LoginReq struct {
	Uid      string `form:"uid"`
	Password string `form:"password"`
}

type LoginResData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRes struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data LoginResData `json:"data"`
}

type VipUpgradeReq struct {
}

type VipUpgradeResData struct {
	VipLevel int `json:"vip_level"`
}

type VipUpgradeRes struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data VipUpgradeResData `json:"data"`
}

type VipResetReq struct {
}

type VipResetRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Ok   bool   `json:"ok"`
}
