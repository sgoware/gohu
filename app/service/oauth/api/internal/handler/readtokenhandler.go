package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gohu/app/service/oauth/api/internal/logic"
	"gohu/app/service/oauth/api/internal/svc"
	"gohu/app/service/oauth/api/internal/types"
)

func ReadTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReadTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewReadTokenLogic(r.Context(), svcCtx)
		resp, err := l.ReadToken(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
