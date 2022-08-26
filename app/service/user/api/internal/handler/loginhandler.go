package handler

import (
	"context"
	"main/app/common/log"
	"main/app/utils/cookie"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
	"main/app/service/user/api/internal/logic"
	"main/app/service/user/api/internal/svc"
	"main/app/service/user/api/internal/types"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		logger := log.GetSugaredLogger()

		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			logger.Errorf("Parse http args failed, \nrequest: \n%v, \nerr: \n%v", r, err)
			return
		}
		logger.Debugf("recv args: %v", req)

		ctx := context.Background()
		if r.RemoteAddr != "" {
			output := strings.Split(r.RemoteAddr, ":")
			ctx = context.WithValue(r.Context(), "lastIp", output[0])
		}
		l := logic.NewLoginLogic(ctx, svcCtx)

		res, err := l.Login(&req)
		if err != nil {
			logger.Errorf("Process logic failed, err: %v", err)
		}

		cookieWriter := cookie.NewCookieWriter(svcCtx.Cookie.Secret,
			cookie.Option{
				Writer:  w,
				Request: r,
				Config:  svcCtx.Cookie.Cookie,
			})

		cookieWriter.Set("x-token", res.Data.AccessToken)
		cookieWriter.Set("refresh-token", res.Data.RefreshToken)

		logger.Info("response: %v", res)
		httpx.WriteJson(w, res.Code, res)
	}
}
