package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Result(w http.ResponseWriter, code int, msg string) {
	httpx.WriteJson(
		w, code, response{
			Code: code,
			Msg:  msg,
		},
	)
}

func ResultWithData(w http.ResponseWriter, code int, msg string, data interface{}) {
	httpx.WriteJson(
		w, code, response{
			Code: code,
			Msg:  msg,
			Data: data,
		},
	)
}
