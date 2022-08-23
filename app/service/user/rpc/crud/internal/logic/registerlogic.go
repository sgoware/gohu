package logic

import (
	"context"
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

func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterRes, error) {
	if len(strings.TrimSpace(in.Uid)) == 0 || len(strings.TrimSpace(in.Password)) == 0 {
		return &crud.RegisterRes{
			Code: http.StatusOK,
			Msg:  "create user failed, err: param err",
		}, nil
	}
	userModel := l.svcCtx.UserModel.User
	_, err := userModel.WithContext(l.ctx).Where(userModel.UID.Eq(gconv.Int64(in.Uid))).First()
	switch err {
	case nil:
		{
			return &crud.RegisterRes{
				Code: http.StatusForbidden,
				Msg:  "create user failed, err: user already exist",
			}, nil
		}
	case gorm.ErrRecordNotFound:
		{
			err := l.svcCtx.UserModel.User.WithContext(l.ctx).Create(&model.User{
				UID:      gconv.Int64(in.Uid),
				Nickname: in.Nickname,
				Password: gmd5.MustEncryptString(in.Password),
			})
			if err != nil {
				return &crud.RegisterRes{
					Code: http.StatusInternalServerError,
					Msg:  "create user failed, err: internal err",
				}, err
			} else {
				return &crud.RegisterRes{
					Code: http.StatusOK,
					Msg:  "create user successfully",
				}, nil
			}
		}
	default:
		return &crud.RegisterRes{
			Code: http.StatusInternalServerError,
			Msg:  "create user failed, err: internal err",
		}, err
	}
}
