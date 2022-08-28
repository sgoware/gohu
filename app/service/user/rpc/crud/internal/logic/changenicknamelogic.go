package logic

import (
	"context"
	"gorm.io/gorm"
	"main/app/common/log"
	"net/http"

	"main/app/service/user/rpc/crud/internal/svc"
	"main/app/service/user/rpc/crud/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeNickNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeNickNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeNickNameLogic {
	return &ChangeNickNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeNickNameLogic) ChangeNickName(in *pb.ChangeNicknameReq) (res *pb.ChangeNicknameRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	userSubjectModel := l.svcCtx.UserModel.UserSubject
	_, err = userSubjectModel.WithContext(l.ctx).Where(userSubjectModel.Nickname.Eq(in.Nickname)).First()
	switch err {
	case nil:
		// 用户名已经存在
		{
			res = &pb.ChangeNicknameRes{
				Code: http.StatusBadRequest,
				Msg:  "nickname already exist",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
	case gorm.ErrRecordNotFound:
		{
			_, err := userSubjectModel.WithContext(l.ctx).
				Where(userSubjectModel.ID.Eq(in.Id)).
				Update(userSubjectModel.Nickname, in.Nickname)
			if err != nil {
				logger.Errorf("change nickname failed, err: %v", err)
				res = &pb.ChangeNicknameRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %v", res.String())
				return res, nil
			}
			res = &pb.ChangeNicknameRes{
				Code: http.StatusOK,
				Msg:  "change nickname successfully",
				Ok:   true,
			}
		}
	default:
		res = &pb.ChangeNicknameRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	logger.Debugf("send message: %v", res.String())
	return res, nil
}
