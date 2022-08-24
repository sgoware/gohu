package logic

import (
	"context"
	"main/app/common/log"
	"main/app/utils/mapping"

	"main/app/service/oauth/rpc/token/enhancer/internal/svc"
	"main/app/service/oauth/rpc/token/enhancer/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserDetailsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDetailsLogic {
	return &GetUserDetailsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserDetailsLogic) GetUserDetails(in *pb.GetUserDetailsReq) (*pb.GetUserDetailsRes, error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	userDetails, err := l.svcCtx.Enhancer.GetUserDetails(in.AccessToken)
	if err != nil {
		logger.Errorf("parse oauth_token failed, err: %v", err)
		return &pb.GetUserDetailsRes{
			Ok:  false,
			Msg: "read oauth_token failed",
		}, err
	}
	logger.Debugf("userDetails: %v", userDetails)
	res := &pb.GetUserDetailsRes{
		Ok:   true,
		Msg:  "get user details successfully",
		Data: &pb.GetUserDetailsRes_Data{UserDetails: &pb.UserDetails{}},
	}
	err = mapping.Struct2Struct(userDetails, res.Data.UserDetails)
	if err != nil {
		logger.Errorf("mapping struct failed, err: %v", err)
		return nil, err
	}

	logger.Debugf("send message: %v", res.String())
	return res, nil
}
