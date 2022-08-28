package logic

import (
	"context"
	"main/app/common/log"
	"main/app/service/oauth/model"
	"main/app/service/oauth/rpc/token/store/internal/svc"
	"main/app/service/oauth/rpc/token/store/pb"
	"net/http"
	"strconv"
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

func (l *StoreTokenLogic) StoreToken(in *pb.StoreTokenReq) (res *pb.StoreTokenRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	if in.UserId == 0 || in.AccessToken == nil {
		res = &pb.StoreTokenRes{
			Code: http.StatusBadRequest,
			Msg:  "invalid param",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	accessTokenString, err := jsonx.MarshalToString(in.AccessToken)
	if err != nil {
		logger.Errorf("marshal access_token to string failed, err: %v", err)
		res = &pb.StoreTokenRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	l.svcCtx.Rdb.Set(l.ctx,
		model.JwtToken+"_"+strconv.FormatInt(in.UserId, 10),
		accessTokenString,
		time.Unix(in.AccessToken.RefreshToken.ExpiresAt, 0).Sub(time.Now()))

	res = &pb.StoreTokenRes{
		Code: http.StatusOK,
		Msg:  "store token successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
