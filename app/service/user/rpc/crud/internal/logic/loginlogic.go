package logic

import (
	"context"
	"encoding/base64"
	"main/app/common/log"
	"main/app/service/user/rpc/crud/crud"
	"main/app/service/user/rpc/crud/internal/svc"
	"main/app/service/user/rpc/crud/pb"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (res *pb.LoginRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())
	// 判断传入参数是否为空
	if len(strings.TrimSpace(in.Uid)) == 0 || len(strings.TrimSpace(in.Password)) == 0 {
		res = &crud.LoginRes{
			Code: http.StatusOK,
			Msg:  "login failed, err: param not fit",
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	// 查找用户
	userModel := l.svcCtx.UserModel.User
	userInfo, err := userModel.WithContext(l.ctx).Where(userModel.UID.Eq(cast.ToInt64(in.Uid))).First()
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		res = &crud.LoginRes{
			Code: http.StatusNotFound,
			Msg:  "login failed, err: uid not exist",
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	default:
		{
			logger.Errorf("database err, err: %v", err)
			return &crud.LoginRes{
				Code: http.StatusInternalServerError,
				Msg:  "login failed, err: internal err",
				Data: nil,
			}, err
		}
	}
	logger.Debugf("userInfo: \n%v", userInfo)
	if userInfo.Password != gmd5.MustEncryptString(in.Password) {
		res = &crud.LoginRes{
			Code: http.StatusUnauthorized,
			Msg:  "login failed, err: wrong password",
			Data: nil,
		}
		logger.Debugf("send message: %v", res.String())

		return res, nil
	}
	encodeAuthString := base64.StdEncoding.EncodeToString([]byte(l.svcCtx.ClientId + ":" + l.svcCtx.ClientSecret + ":" + cast.ToString(userInfo.UID)))
	logger.Debugf("encodeAuthString: %v", encodeAuthString)
	basicAuthString := "Basic " + encodeAuthString
	logger.Debugf("basicAuthString: %v", basicAuthString)
	res = &crud.LoginRes{
		Code: http.StatusOK,
		Msg:  "get auth_token successfully",
		Data: &crud.LoginRes_Data{AuthToken: basicAuthString},
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
