package logic

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"main/app/common/log"
	"net/http"

	"main/app/service/user/rpc/info/internal/svc"
	"main/app/service/user/rpc/info/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPersonalInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPersonalInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPersonalInfoLogic {
	return &GetPersonalInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPersonalInfoLogic) GetPersonalInfo(in *pb.GetPersonalInfoReq) (res *pb.GetPersonalInfoRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	userSubjectBytes, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("user_subject_%d", in.UserId)).Bytes()
	if err == nil {
		rpcResData := &pb.GetPersonalInfoRes_Data{}
		err = proto.Unmarshal(userSubjectBytes, rpcResData)
		if err != nil {
			logger.Errorf("unmarshal [rpcResData] failed, err: %v", err)
		} else {
			res = &pb.GetPersonalInfoRes{
				Code: http.StatusOK,
				Msg:  "get personal info successfully",
				Ok:   true,
				Data: rpcResData,
			}
			logger.Debugf("send message: %v", err)
			return res, nil
		}
	}
	logger.Errorf("get [user_subject] cache failed, err: %v", err)

	userSubjectModel := l.svcCtx.UserModel.UserSubject

	userSubject, err := userSubjectModel.WithContext(l.ctx).
		Where(userSubjectModel.ID.Eq(in.UserId)).
		First()
	switch err {
	case gorm.ErrRecordNotFound:
		res = &pb.GetPersonalInfoRes{
			Code: http.StatusForbidden,
			Msg:  "user not found",
			Ok:   false,
			Data: nil,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil

	case nil:

	default:
		logger.Debugf("get personal info failed, err: mysql err, %v", err)
		res = &pb.GetPersonalInfoRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
			Data: nil,
		}
		return res, nil
	}

	res = &pb.GetPersonalInfoRes{
		Code: http.StatusOK,
		Msg:  "get personal info successfully",
		Ok:   true,
		Data: &pb.GetPersonalInfoRes_Data{
			Username: userSubject.Username,
			Nickname: userSubject.Nickname,
			Email:    userSubject.Email,
			Phone:    userSubject.Phone,
			Vip:      userSubject.Vip,
		},
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
