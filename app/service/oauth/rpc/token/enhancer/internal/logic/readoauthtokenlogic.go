package logic

import (
	"context"
	"main/app/common/log"
	"main/app/service/oauth/rpc/token/enhancer/internal/svc"
	"main/app/service/oauth/rpc/token/enhancer/pb"
	"main/app/utils/mapping"

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

func (l *ReadOauthTokenLogic) ReadOauthToken(in *pb.ReadTokenReq) (*pb.ReadTokenRes, error) {
	logger := log.GetSugaredLogger()

	oauthToken, _, err := l.svcCtx.Enhancer.ParseToken(in.OauthToken)
	if err != nil {
		logger.Errorf("parse oauth_token failed, err: %v", err)
		return &pb.ReadTokenRes{
			Ok:  false,
			Msg: "read oauth_token failed",
		}, err
	}
	oauthTokenRes := &pb.ReadTokenRes{
		Ok:   true,
		Msg:  "read oauth_token successfully",
		Data: &pb.ReadTokenRes_Data{AccessToken: &pb.OAuth2Token{}},
	}
	err = mapping.Struct2Struct(oauthToken, oauthTokenRes.Data.AccessToken)
	if err != nil {
		return nil, err
	}
	return oauthTokenRes, nil
}
