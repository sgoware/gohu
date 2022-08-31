package logic

import (
	"context"
	"main/app/common/log"
	"main/app/service/user/dao/model"
	"main/app/service/user/rpc/crud/internal/svc"
	"main/app/service/user/rpc/crud/pb"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSubscriptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSubscriptionLogic {
	return &CreateSubscriptionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateSubscriptionLogic) CreateSubscription(in *pb.CreateSubscriptionReq) (res *pb.CreateSubscriptionRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	userSubscriptionModel := l.svcCtx.UserModel.UserSubscription

	err = userSubscriptionModel.WithContext(l.ctx).
		Create(&model.UserSubscription{
			UserID:  in.UserId,
			ObjType: in.ObjType,
			ObjID:   in.ObjId,
		})
	if err != nil {
		logger.Errorf("create subscription failed, err: mysql err, %v", err)
		res = &pb.CreateSubscriptionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	res = &pb.CreateSubscriptionRes{
		Code: http.StatusOK,
		Msg:  "create subscription successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
