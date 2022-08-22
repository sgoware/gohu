package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gohu/app/service/oauth/api/internal/logic"
	"gohu/app/service/oauth/api/internal/svc"
	"gohu/app/service/oauth/api/internal/types"
)

func RefreshTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetTokenByRefreshTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewRefreshTokenLogic(r.Context(), svcCtx)
		resp, err := l.RefreshToken(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
