package logic

import (
	"context"
	"main/app/common/log"
	"main/app/service/user/dao/model"
	"main/app/service/user/rpc/crud/crud"
	"main/app/service/user/rpc/crud/internal/svc"
	"main/app/service/user/rpc/crud/pb"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/util/gconv"
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

	if len(strings.TrimSpace(in.Uid)) == 0 || len(strings.TrimSpace(in.Password)) == 0 {
		res = &crud.RegisterRes{
			Code: http.StatusOK,
			Msg:  "create user failed, err: param err",
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
	userModel := l.svcCtx.UserModel.User
	_, err = userModel.WithContext(l.ctx).Where(userModel.UID.Eq(gconv.Int64(in.Uid))).First()
	switch err {
	case nil:
		{
			res = &crud.RegisterRes{
				Code: http.StatusForbidden,
				Msg:  "create user failed, err: user already exist",
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
	case gorm.ErrRecordNotFound:
		{
			err := l.svcCtx.UserModel.User.WithContext(l.ctx).Create(&model.User{
				UID:      gconv.Int64(in.Uid),
				Nickname: in.Nickname,
				Password: gmd5.MustEncryptString(in.Password),
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
	default:
		res = &crud.RegisterRes{
			Code: http.StatusInternalServerError,
			Msg:  "create user failed, err: internal err",
		}
		logger.Debugf("send message: %v", res.String())
		return res, err
	}
}
