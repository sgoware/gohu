// Code generated by goctl. DO NOT EDIT.
package types

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
	Id int64 `path:"question_id"`
}

type GetQuestionRes struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Ok   bool   `json:"ok"`
	Data string `json:"data"`
}

type GetAnswerReq struct {
	Id int64 `path:"id"`
}

type GetAnswerRes struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Ok   bool   `json:"ok"`
	Data string `json:"data"`
}
