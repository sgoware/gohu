package logic

import (
	"context"
	"main/app/common/log"
	"net/http"

	"main/app/service/user/rpc/vip/internal/svc"
	"main/app/service/user/rpc/vip/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetLogic {
	return &ResetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetLogic) Reset(in *pb.ResetReq) (res *pb.ResetRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	userModel := l.svcCtx.UserModel
	_, err = userModel.WithContext(l.ctx).User.Select(userModel.User.ID, userModel.User.Vip).
		Where(userModel.User.ID.Eq(in.Id)).Update(userModel.User.Vip, 0)
	if err != nil {
		logger.Errorf("reset vip failed, err: %v", err)
		return &pb.ResetRes{
			Code: http.StatusOK,
			Msg:  "internal err",
			Ok:   false,
		}, nil
	}

	res = &pb.ResetRes{
		Code: http.StatusOK,
		Msg:  "reset vip level successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
