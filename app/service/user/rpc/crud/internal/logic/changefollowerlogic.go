package logic

import (
	"context"
	"fmt"
	"main/app/common/log"
	"net/http"

	"main/app/service/user/rpc/crud/internal/svc"
	"main/app/service/user/rpc/crud/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeFollowerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeFollowerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeFollowerLogic {
	return &ChangeFollowerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeFollowerLogic) ChangeFollower(in *pb.ChangeFollowerReq) (res *pb.ChangeFollowerRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	l.svcCtx.Rdb.SAdd(l.ctx,
		"user_follower",
		fmt.Sprintf("%d:%d", in.UserId))

	userSubjectModel := l.svcCtx.UserModel.UserSubject

	userSubject, err := userSubjectModel.WithContext(l.ctx).
		Select(userSubjectModel.ID, userSubjectModel.Follower).
		Where(userSubjectModel.ID.Eq(in.UserId)).First()
	if err != nil {
		logger.Errorf("change follower failed, err: mysql err, %v", err)
		res = &pb.ChangeFollowerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	switch in.Action {
	case 1:
		// 增加
		_, err = userSubjectModel.WithContext(l.ctx).
			Where(userSubjectModel.ID.Eq(in.UserId)).
			Update(userSubjectModel.Follower, userSubject.Follower+1)
		if err != nil {
			logger.Errorf("change follower failed, err: mysql err, %v", err)
			res = &pb.ChangeFollowerRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
	case 2:
		// 删除
		_, err = userSubjectModel.WithContext(l.ctx).
			Where(userSubjectModel.ID.Eq(in.UserId)).
			Update(userSubjectModel.Follower, userSubject.Follower-1)
		if err != nil {
			logger.Errorf("change follower failed, err: mysql err, %v", err)
			res = &pb.ChangeFollowerRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
	}

	res = &pb.ChangeFollowerRes{
		Code: http.StatusOK,
		Msg:  "change follower successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
