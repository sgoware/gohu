package logic

import (
	"context"
	"encoding/base64"
	"main/app/service/user/rpc/crud/crud"
	"main/app/service/user/rpc/crud/internal/svc"
	"main/app/service/user/rpc/crud/pb"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/util/gconv"
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

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginRes, error) {
	// 判断传入参数是否为空
	if len(strings.TrimSpace(in.Uid)) == 0 || len(strings.TrimSpace(in.Password)) == 0 {
		return &crud.LoginRes{
			Code: http.StatusOK,
			Msg:  "login failed, err: param not fit",
		}, nil
	}

	// 查找用户
	userModel := l.svcCtx.UserModel.User
	userInfo, err := userModel.WithContext(l.ctx).Where(userModel.UID.Eq(gconv.Int64(in.Uid))).First()
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		return &crud.LoginRes{
			Code: http.StatusNotFound,
			Msg:  "login failed, err: uid not exist",
		}, nil
	default:
		{
			logx.Errorf("database err, err: %v", err)
			return &crud.LoginRes{
				Code: http.StatusInternalServerError,
				Msg:  "login failed, err: internal err",
				Data: nil,
			}, err
		}
	}
	if userInfo.Password != gmd5.MustEncryptString(in.Password) {
		return &crud.LoginRes{
			Code: http.StatusUnauthorized,
			Msg:  "login failed, err: wrong password",
			Data: nil,
		}, nil
	}
	encodeAuthString := base64.StdEncoding.EncodeToString([]byte(l.svcCtx.ClientId + ":" + l.svcCtx.ClientSecret + ":" + cast.ToString(userInfo.UID)))
	BasicAuthString := "Basic " + encodeAuthString
	return &crud.LoginRes{
		Code: http.StatusOK,
		Msg:  "get auth_token successfully",
		Data: &crud.LoginRes_Data{AuthToken: BasicAuthString},
	}, nil
}
