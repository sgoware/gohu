package logic

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"golang.org/x/crypto/sha3"
	"main/app/common/log"
	"main/app/service/user/rpc/crud/crud"
	"main/app/service/user/rpc/crud/internal/svc"
	"main/app/service/user/rpc/crud/pb"
	"net/http"
	"strings"

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
	if len(strings.TrimSpace(in.Username)) == 0 || len(strings.TrimSpace(in.Password)) == 0 {
		res = &crud.LoginRes{
			Code: http.StatusBadRequest,
			Msg:  "param not fit",
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	// 在数据库中查找用户
	userSubjectModel := l.svcCtx.UserModel.UserSubject
	userInfo, err := userSubjectModel.WithContext(l.ctx).Where(userSubjectModel.Username.Eq(in.Username)).First()
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		res = &crud.LoginRes{
			Code: http.StatusNotFound,
			Msg:  "uid not exist",
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	default:
		{
			logger.Errorf("database err, err: %v", err)
			res = &crud.LoginRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Data: nil,
			}
			logger.Debugf("send message: %v", res.String())
			return res, err
		}
	}
	logger.Debugf("userInfo: \n%v", userInfo)

	// 验证密码
	d := sha3.Sum224([]byte(in.Password))
	encryptedPassword := hex.EncodeToString(d[:])
	if userInfo.Password != encryptedPassword {
		res = &crud.LoginRes{
			Code: http.StatusUnauthorized,
			Msg:  "wrong password",
			Data: nil,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	// 更新最近登录 ip
	_, err = userSubjectModel.WithContext(l.ctx).
		Where(userSubjectModel.Username.Eq(in.Username)).
		Update(userSubjectModel.LastIP, in.LastIp)
	if err != nil {
		logger.Errorf("database err, err: %v", err)
		res = &crud.LoginRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Data: nil,
		}
		return res, err
	}

	// 生成 oauth 服务器的认证头
	encodeAuthString := base64.StdEncoding.EncodeToString([]byte(l.svcCtx.ClientId + ":" + l.svcCtx.ClientSecret + ":" + cast.ToString(userInfo.ID)))
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
