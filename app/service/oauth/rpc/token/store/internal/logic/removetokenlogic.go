package logic

import (
	"context"
	"gohu/app/common/log"
	"gohu/app/service/oauth/model"
	"gohu/app/service/oauth/rpc/token/store/internal/svc"
	"gohu/app/service/oauth/rpc/token/store/pb"

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

func (l *RemoveTokenLogic) RemoveToken(in *pb.RemoveTokenReq) (*pb.RemoveTokenRes, error) {
	logger := log.GetSugaredLogger()

	if in.UserId == " " {
		logger.Error("remove token failed, err: %v", model.ErrInvalidTokenRequest)
		return &pb.RemoveTokenRes{
			Ok:  false,
			Msg: "remove token failed, err: invalid token request",
		}, nil
	}

	l.svcCtx.Rdb.Set(l.ctx, model.JwtToken+" "+in.UserId, 1, 0)

	return &pb.RemoveTokenRes{
		Ok:  true,
		Msg: "remove token successfully",
	}, nil
}
