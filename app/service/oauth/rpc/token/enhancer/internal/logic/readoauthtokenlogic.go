package logic

import (
	"context"
	"fmt"
	"main/app/common/log"
	"main/app/service/oauth/rpc/token/enhancer/internal/svc"
	"main/app/service/oauth/rpc/token/enhancer/pb"
	"main/app/utils/mapping"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadOauthTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReadOauthTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadOauthTokenLogic {
	return &ReadOauthTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReadOauthTokenLogic) ReadOauthToken(in *pb.ReadTokenReq) (res *pb.ReadTokenRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	oauthToken, _, err := l.svcCtx.Enhancer.ParseToken(in.OauthToken)
	if err != nil {
		res = &pb.ReadTokenRes{
			Code: http.StatusOK,
			Msg:  fmt.Sprintf("parse oauth token failed, %v", err),
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, err
	}
	logger.Debugf("oauthToken: %v", oauthToken)
	res = &pb.ReadTokenRes{
		Code: http.StatusOK,
		Msg:  "read oauth token successfully",
		Ok:   true,
		Data: &pb.ReadTokenRes_Data{AccessToken: &pb.OAuth2Token{}},
	}

	err = mapping.Struct2Struct(oauthToken, res.Data.AccessToken)
	if err != nil {
		res = &pb.ReadTokenRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, err
	}

	logger.Debugf("send message: %v", res.String())
	return res, nil
}
