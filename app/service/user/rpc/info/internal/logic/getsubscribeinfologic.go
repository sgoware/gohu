package logic

import (
	"context"
	"main/app/common/log"
	"net/http"

	"main/app/service/user/rpc/info/internal/svc"
	"main/app/service/user/rpc/info/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSubscribeInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSubscribeInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSubscribeInfoLogic {
	return &GetSubscribeInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSubscribeInfoLogic) GetSubscribeInfo(in *pb.GetSubscribeInfoReq) (res *pb.GetSubscribeInfoRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	userSubscribeModel := l.svcCtx.UserModel.UserSubscribe

	userSubscribes, err := userSubscribeModel.WithContext(l.ctx).
		Where(userSubscribeModel.UserID.Eq(in.UserId), userSubscribeModel.ObjType.Eq(in.ObjType)).Find()
	if err != nil {
		logger.Errorf("get subscribe info failed, err: mysql err, %v", err)
		res = &pb.GetSubscribeInfoRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
			Data: nil,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	res = &pb.GetSubscribeInfoRes{
		Code: http.StatusOK,
		Msg:  "get subscribe info successfully",
		Ok:   true,
		Data: &pb.GetSubscribeInfoRes_Data{Ids: make([]int64, 0)},
	}
	for _, userSubscribe := range userSubscribes {
		res.Data.Ids = append(res.Data.Ids, userSubscribe.ObjID)
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
