package logic

import (
	"context"
	"main/app/common/log"
	"net/http"

	"main/app/service/user/rpc/crud/internal/svc"
	"main/app/service/user/rpc/crud/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSubscriptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSubscriptionLogic {
	return &DeleteSubscriptionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteSubscriptionLogic) DeleteSubscription(in *pb.DeleteSubscriptionReq) (res *pb.DeleteSubscriptionRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	userSubscriptionModel := l.svcCtx.UserModel.UserSubscription

	_, err = userSubscriptionModel.WithContext(l.ctx).
		Where(userSubscriptionModel.ID.Eq(in.SubscriptionId)).
		Delete()
	if err != nil {
		logger.Errorf("delete subscription failed, err: mysql err, %v", err)
		res = &pb.DeleteSubscriptionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	res = &pb.DeleteSubscriptionRes{
		Code: http.StatusOK,
		Msg:  "delete subscription successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
