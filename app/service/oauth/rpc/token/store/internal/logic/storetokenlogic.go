package logic

import (
	"context"
	"main/app/common/log"
	"main/app/service/oauth/model"
	"main/app/service/oauth/rpc/token/store/internal/svc"
	"main/app/service/oauth/rpc/token/store/pb"
	"time"

	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
)

type StoreTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStoreTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StoreTokenLogic {
	return &StoreTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StoreTokenLogic) StoreToken(in *pb.StoreTokenReq) (*pb.StoreTokenRes, error) {
	logger := log.GetSugaredLogger()

	if in.UserId == " " || in.AccessToken == nil {
		logger.Errorf("store token failed, err: %v", model.ErrInvalidTokenRequest)
		return &pb.StoreTokenRes{
			Ok:  false,
			Msg: "store token failed, err: invalid token request",
		}, nil
	}

	accessTokenString, err := jsonx.MarshalToString(in.AccessToken)
	if err != nil {
		logger.Errorf("marshal access_token to string failed, err: %v", err)
		return &pb.StoreTokenRes{
			Ok:  false,
			Msg: "marshal access_token to string failed",
		}, nil
	}
	logger.Debugf("%v", accessTokenString)
	l.svcCtx.Rdb.Set(l.ctx, model.JwtToken+"_"+in.UserId, accessTokenString, time.Unix(in.AccessToken.ExpiresAt, 0).Sub(time.Now()))

	return &pb.StoreTokenRes{
		Ok:  true,
		Msg: "store token successfully",
	}, nil
}
