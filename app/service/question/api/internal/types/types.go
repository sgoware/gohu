// Code generated by goctl. DO NOT EDIT.
package types

type QuestionSubject struct {
	Id          int64  `json:"id"`
	UserId      int64  `json:"user_id"`
	IpLoc       string `json:"ip_loc"`
	Title       string `json:"title"`
	Topic       string `json:"topic"`
	Tag         string `json:"tag"`
	SubCount    int32  `json:"sub_count"`
	AnswerCount int32  `json:"answer_count"`
	ViewCount   int64  `json:"view_count"`
	State       int32  `json:"state"`
	Attr        int32  `json:"attr"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
}

type QuestionContent struct {
	QuestionId int64  `json:"question_id"`
	Content    string `json:"content"`
	Meta       string `json:"meta"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

type AnswerIndex struct {
	Id           int64  `json:"id"`
	QuestionId   int64  `json:"question_id"`
	UserId       int64  `json:"user_id"`
	IpLoc        string `json:"ip_loc"`
	ApproveCount int32  `json:"approve_count"`
	LikeCount    int32  `json:"like_count"`
	CollectCount int32  `json:"collect_count"`
	State        int32  `json:"state"`
	Attrs        int32  `json:"attrs"`
	CreateTime   string `json:"create_time"`
	UpdateTime   string `json:"update_time"`
}

type AnswerContent struct {
	AnswerId   int64  `json:"answer_id"`
	Content    string `json:"content"`
	Meta       string `json:"meta"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

type CrudReq struct {
	Object string `form:"object"`
	Action string `form:"action"`
	Data   string `form:"data"`
}

type CrudRes struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Ok   bool   `json:"ok"`
}

type GetQuestionReq struct {
	QuestionId int64 `path:"question_id"`
}

type GetQuestionResData struct {
	QuestionSubject QuestionSubject `json:"question_subject"`
	QuestionContent QuestionContent `json:"question_content"`
}

type GetQuestionRes struct {
	Code int32              `json:"code"`
	Msg  string             `json:"msg"`
	Ok   bool               `json:"ok"`
	Data GetQuestionResData `json:"data"`
}

type GetAnswerReq struct {
	AnswerId int64 `path:"answer_id"`
}

type GetAnswerResData struct {
	AnswerIndex   AnswerIndex   `json:"answer_index"`
	AnswerContent AnswerContent `json:"answer_content"`
}

type GetAnswerRes struct {
	Code int32            `json:"code"`
	Msg  string           `json:"msg"`
	Ok   bool             `json:"ok"`
	Data GetAnswerResData `json:"data"`
}
