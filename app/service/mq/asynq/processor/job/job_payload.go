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

type MsgCrudQuestionSubjectRecordPayload struct {
	Action      int32     `json:"action"` // 1:创建 2:更新 3:删除
	Id          int64     `json:"id"`
	UserId      int64     `json:"user_id"`
	IpLoc       string    `json:"ip_loc"`
	Title       string    `json:"title"`
	Topic       string    `json:"topic"`
	Tag         string    `json:"tag"`
	SubCount    int32     `json:"sub_count"`
	AnswerCount int32     `json:"answer_count"`
	ViewCount   int64     `json:"view_count"`
	State       int32     `json:"state"`
	Attrs       int32     `json:"attrs"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

type MsgCrudQuestionContentRecordPayload struct {
	Action     int32     `json:"action"`
	QuestionId int64     `json:"question_id"`
	Content    string    `json:"content"`
	Meta       string    `json:"meta"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

type MsgCrudCommentSubjectPayload struct {
	Action     int32     `json:"action"` // 1:创建 2:更新 3:删除
	Id         int64     `json:"id"`
	ObjType    int32     `json:"obj_type"`
	ObjId      int64     `json:"obj_id"`
	Count      int32     `json:"count"`
	RootCount  int32     `json:"root_count"`
	State      int32     `json:"state"`
	Attrs      int32     `json:"attrs"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

type ScheduleUpdateCommentRecordPayload struct {
}
