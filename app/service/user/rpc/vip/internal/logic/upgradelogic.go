package logic

import (
	"context"
	"github.com/spf13/cast"
	"main/app/common/log"
	"net/http"

	"main/app/service/user/rpc/vip/internal/svc"
	"main/app/service/user/rpc/vip/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpgradeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpgradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpgradeLogic {
	return &UpgradeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpgradeLogic) Upgrade(in *pb.UpgradeReq) (res *pb.UpgradeRes, err error) {
	// todo: add your logic here and delete this line
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	userModel := l.svcCtx.UserModel
	userInfo, _ := userModel.WithContext(l.ctx).User.Select(userModel.User.ID, userModel.User.Vip).
		Where(userModel.User.ID.Eq(cast.ToInt64(in.Id))).First()
	if userInfo.Vip < 9 {
		userModel.WithContext(l.ctx).User.Select(userModel.User.ID, userModel.User.Vip).
			Where(userModel.User.ID.Eq(userInfo.ID)).
			Update(userModel.User.Vip, userInfo.Vip+1)
		res = &pb.UpgradeRes{
			Code: http.StatusOK,
			Msg:  "vip upgrade successfully",
			Data: &pb.UpgradeRes_Data{VipLevel: userInfo.Vip + 1},
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	} else {
		res = &pb.UpgradeRes{
			Code: http.StatusBadRequest,
			Msg:  "vip level is already the highest",
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
}
