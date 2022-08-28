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

type CheckTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckTokenLogic {
	return &CheckTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckTokenLogic) CheckToken(in *pb.CheckTokenReq) (res *pb.CheckTokenRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	if in.UserId == " " {
		res = &pb.CheckTokenRes{
			Code:    http.StatusBadRequest,
			Msg:     model.ErrInvalidUserId.Error(),
			Ok:      false,
			IsExist: false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	_, err = l.svcCtx.Rdb.Get(l.ctx, model.JwtToken+" "+in.UserId).Result()
	if err != nil {
		res = &pb.CheckTokenRes{
			Code:    http.StatusOK,
			Msg:     "token is not exist",
			Ok:      true,
			IsExist: false,
		}
		return res, nil
	}
	res = &pb.CheckTokenRes{
		Code:    http.StatusOK,
		Msg:     "token is exist",
		Ok:      true,
		IsExist: true,
	}
	return res, nil
}
