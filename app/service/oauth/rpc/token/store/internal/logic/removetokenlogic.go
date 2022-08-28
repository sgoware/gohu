package logic

import (
	"context"
	"main/app/common/log"
	"main/app/service/oauth/model"
	"main/app/service/oauth/rpc/token/store/internal/svc"
	"main/app/service/oauth/rpc/token/store/pb"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveTokenLogic {
	return &RemoveTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveTokenLogic) RemoveToken(in *pb.RemoveTokenReq) (res *pb.RemoveTokenRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	if in.UserId == " " {
		res = &pb.RemoveTokenRes{
			Code: http.StatusBadRequest,
			Msg:  model.ErrInvalidUserId.Error(),
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	l.svcCtx.Rdb.Set(l.ctx, model.JwtToken+" "+in.UserId, 1, 0)

	res = &pb.RemoveTokenRes{
		Code: http.StatusOK,
		Msg:  "remove token successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res)
	return res, nil
}
