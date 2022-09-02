package job

import "time"

type UserSubjectPayload struct {
	Id         int64     `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Nickname   string    `json:"nickname"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	LastIp     string    `json:"last_ip"`
	Vip        int32     `json:"vip"`
	Follower   int32     `json:"follower"`
	State      int32     `json:"state"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

type MsgCreateUserSubjectPayload struct {
	Username   string    `json:"username"`
	Nickname   string    `json:"nickname"`
	Password   string    `json:"password"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

type MsgAddUserSubjectCachePayload struct {
	Id       int64 `json:"id"`
	Vip      int32 `json:"vip"`
	Follower int32 `json:"follower"`
}

type UserCollectPayload struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"userId"`
	CollectType int32     `json:"collectType"`
	ObjType     int32     `json:"objType"`
	ObjId       int64     `json:"objId"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

type ScheduleUpdateQuestionRecordPayload struct {
}

type ScheduleUpdateCommentRecordPayload struct {
}
