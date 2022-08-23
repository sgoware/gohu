package logic

import (
	"context"
	"main/app/common/log"
	"main/app/service/oauth/model"
	"main/app/service/oauth/rpc/token/store/internal/svc"
	"main/app/service/oauth/rpc/token/store/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
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

func (l *CheckTokenLogic) CheckToken(in *pb.CheckTokenReq) (*pb.CheckTokenRes, error) {
	logger := log.GetSugaredLogger()

	if in.UserId == " " {
		logger.Errorf("check token failed, err: %v", model.ErrInvalidTokenRequest)
		return &pb.CheckTokenRes{
			Ok:      false,
			Msg:     "check token failed, err: invalid token request",
			IsExist: false,
		}, nil
	}

	_, err := l.svcCtx.Rdb.Get(l.ctx, model.JwtToken+" "+in.UserId).Result()
	if err == nil {
		return &pb.CheckTokenRes{
			Ok:      true,
			Msg:     "check token successfully",
			IsExist: true,
		}, nil
	} else {
		if err != redis.ErrEmptyKey {
			return &pb.CheckTokenRes{
				Ok:      false,
				Msg:     "check token failed, err: redis err",
				IsExist: false,
			}, nil
		} else {
			return &pb.CheckTokenRes{
				Ok:      true,
				Msg:     "check token successfully",
				IsExist: false,
			}, nil
		}
	}
}
