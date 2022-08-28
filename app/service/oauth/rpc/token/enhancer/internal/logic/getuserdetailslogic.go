package logic

import (
	"context"
	"fmt"
	"main/app/common/log"
	"main/app/utils/mapping"
	"net/http"

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

func (l *GetUserDetailsLogic) GetUserDetails(in *pb.GetUserDetailsReq) (res *pb.GetUserDetailsRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	userDetails, err := l.svcCtx.Enhancer.GetUserDetails(in.AccessToken)
	if err != nil {
		res = &pb.GetUserDetailsRes{
			Code: http.StatusOK,
			Msg:  fmt.Sprintf("parse oauth token failed, %v", err),
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, err
	}
	res = &pb.GetUserDetailsRes{
		Ok:   true,
		Msg:  "get user details from oauth token successfully",
		Data: &pb.GetUserDetailsRes_Data{UserDetails: &pb.UserDetails{}},
	}
	err = mapping.Struct2Struct(userDetails, res.Data.UserDetails)
	if err != nil {
		res = &pb.GetUserDetailsRes{
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
