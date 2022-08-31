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

type CreateCollectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCollectionLogic {
	return &CreateCollectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCollectionLogic) CreateCollection(in *pb.CreateCollectionReq) (res *pb.CreateCollectionRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	userCollectionModel := l.svcCtx.UserModel.UserCollection

	err = userCollectionModel.WithContext(l.ctx).
		Create(&model.UserCollection{
			UserID:  in.UserId,
			ObjType: in.ObjType,
			ObjID:   in.ObjId,
		})
	if err != nil {
		logger.Errorf("create collection failed, err: mysql err, %v", err)
		res = &pb.CreateCollectionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	res = &pb.CreateCollectionRes{
		Code: http.StatusOK,
		Msg:  "create collection successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
