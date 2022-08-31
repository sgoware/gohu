// Code generated by goctl. DO NOT EDIT.
package types

type RegisterReq struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type RegisterRes struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Ok   bool   `json:"ok"`
}

type LoginReq struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type LoginResData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRes struct {
	Code int32        `json:"code"`
	Msg  string       `json:"msg"`
	Ok   bool         `json:"ok"`
	Data LoginResData `json:"data"`
}

type ChangeNicknameReq struct {
	Nickname string `form:"nickname"`
}

type ChangeNicknameRes struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Ok   bool   `json:"ok"`
}

type CreateCollectionReq struct {
	ObjType int32 `form:"obj_type"`
	ObjId   int64 `form:"obj_id"`
}

type CreateCollectionRes struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Ok   bool   `json:"ok"`
}

type DeleteCollectionReq struct {
	CollectionId int64 `path:"collection_id"`
}

type DeleteCollectionRes struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Ok   bool   `json:"ok"`
}

type CreateSubscriptionReq struct {
	ObjType int32 `form:"obj_type"`
	ObjId   int64 `form:"obj_id"`
}

type CreateSubscriptionRes struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Ok   bool   `json:"ok"`
}

type DeleteSubscriptionReq struct {
	SubscriptionId int64 `path:"subscription_id"`
}

type DeleteSubscriptionRes struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Ok   bool   `json:"ok"`
}

type GetObjInfoReq struct {
	ObjType int32 `path:"obj_type"`
}

type GetObjInfoResData struct {
	Ids []int64 `json:"ids"`
}

type GetObjInfoRes struct {
	Code int32             `json:"code"`
	Msg  string            `json:"msg"`
	Ok   bool              `json:"ok"`
	Data GetObjInfoResData `json:"data"`
}

type GetPersonalInfoReq struct {
}

type GetPersonalInfoResData struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Vip      int32  `json:"vip"`
}

type GetPersonalInfoRes struct {
	Code int32                  `json:"code"`
	Msg  string                 `json:"msg"`
	Ok   bool                   `json:"ok"`
	Data GetPersonalInfoResData `json:"data"`
}

type GetCollectionInfoReq struct {
}

type GetCollectionInfoResData struct {
	ObjType []int32 `json:"obj_type"`
	ObjId   []int64 `json:"obj_id"`
}

type GetCollectionInfoRes struct {
	Code int32                    `json:"code"`
	Msg  string                   `json:"msg"`
	Ok   bool                     `json:"ok"`
	Data GetCollectionInfoResData `json:"data"`
}

type GetNotificationInfoReq struct {
	MessageType int32 `path:"message_type"`
}

type GetNotificationInfoResData struct {
	MessageIds []int64 `json:"message_ids"`
}

type GetNotificationInfoRes struct {
	Code int32                      `json:"code"`
	Msg  string                     `json:"msg"`
	Ok   bool                       `json:"ok"`
	Data GetNotificationInfoResData `json:"data"`
}

type GetSubscribeInfoReq struct {
	ObjType int32 `path:"obj_type"`
}

type GetSubscribeInfoResData struct {
	Ids []int64 `json:"ids"`
}

type GetSubscribeInfoRes struct {
	Code int32                   `json:"code"`
	Msg  string                  `json:"msg"`
	Ok   bool                    `json:"ok"`
	Data GetSubscribeInfoResData `json:"data"`
}

type VipUpgradeReq struct {
}

type VipUpgradeResData struct {
	VipLevel int32 `json:"vip_level"`
}

type VipUpgradeRes struct {
	Code int32             `json:"code"`
	Msg  string            `json:"msg"`
	Ok   bool              `json:"ok"`
	Data VipUpgradeResData `json:"data"`
}

type VipResetReq struct {
}

type VipResetRes struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Ok   bool   `json:"ok"`
}
