package logic

import (
	"context"
	"encoding/hex"
	"golang.org/x/crypto/sha3"
	"main/app/common/log"
	"main/app/service/user/dao/model"
	"main/app/service/user/rpc/crud/crud"
	"main/app/service/user/rpc/crud/internal/svc"
	"main/app/service/user/rpc/crud/pb"
	"main/app/utils/uuid"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterReq) (res *pb.RegisterRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	if len(strings.TrimSpace(in.Username)) == 0 || len(strings.TrimSpace(in.Password)) == 0 {
		res = &crud.RegisterRes{
			Code: http.StatusOK,
			Msg:  "create user failed, err: param err",
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
	userModel := l.svcCtx.UserModel.User
	_, err = userModel.WithContext(l.ctx).Where(userModel.Username.Eq(in.Username)).First()
	switch err {
	case nil:
		// 用户已经存在的情况
		{
			res = &crud.RegisterRes{
				Code: http.StatusForbidden,
				Msg:  "create user failed, err: user already exist",
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
	case gorm.ErrRecordNotFound:
		// 用户不存在的情况(可以创建用户)
		{
			// 密码使用sha3哈希然后存储
			d := sha3.Sum224([]byte(in.Password))
			encryptedPassword := hex.EncodeToString(d[:])

			// 生成默认昵称
			defaultNickname := uuid.NewRandomString(in.Username, "username", 10)

			err := l.svcCtx.UserModel.User.WithContext(l.ctx).Create(&model.User{
				Username: in.Username,
				Password: encryptedPassword,
				Nickname: "gohu_" + defaultNickname,
			})

			if err != nil {
				res = &crud.RegisterRes{
					Code: http.StatusInternalServerError,
					Msg:  "create user failed, err: internal err",
				}
				logger.Debugf("send message: %v", res.String())
				return res, err
			} else {
				res = &crud.RegisterRes{
					Code: http.StatusOK,
					Msg:  "create user successfully",
				}
				logger.Debugf("send message: %v", res.String())
				return res, nil
			}
		}
		// 数据库查询失败的情况
	default:
		res = &crud.RegisterRes{
			Code: http.StatusInternalServerError,
			Msg:  "create user failed, err: internal err",
		}
		logger.Debugf("send message: %v", res.String())
		return res, err
	}
}
