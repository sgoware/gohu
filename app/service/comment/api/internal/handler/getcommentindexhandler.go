package handler

import (
	"main/app/common/log"
	"main/app/service/comment/api/internal/logic"
	"main/app/service/comment/api/internal/svc"
	"main/app/service/comment/api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCommentIndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCommentIndexReq
		logger := log.GetSugaredLogger()

		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			logger.Errorf("Parse http args failed, \nrequest: \n%v, \nerr: \n%v", r, err)
			return
		}
		logger.Debugf("recv args: %v", req)

		l := logic.NewGetCommentIndexLogic(r.Context(), svcCtx)

		res, err := l.GetCommentIndex(&req)
		if err != nil {
			logger.Errorf("Process logic failed, err: %v", err)
		}

		logger.Info("response: %v", res)
		httpx.WriteJson(w, int(res.Code), res)
	}
}
