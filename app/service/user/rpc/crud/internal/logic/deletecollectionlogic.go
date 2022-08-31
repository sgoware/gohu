package logic

import (
	"context"
	"main/app/common/log"
	"net/http"

	"main/app/service/user/rpc/crud/internal/svc"
	"main/app/service/user/rpc/crud/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCollectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCollectionLogic {
	return &DeleteCollectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCollectionLogic) DeleteCollection(in *pb.DeleteCollectionReq) (res *pb.DeleteCollectionRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	userCollectionModel := l.svcCtx.UserModel.UserCollection

	_, err = userCollectionModel.WithContext(l.ctx).
		Where(userCollectionModel.ID.Eq(in.CollectionId)).
		Delete()
	if err != nil {
		logger.Errorf("delete collection failed, err: mysql err, %v", err)
		res = &pb.DeleteCollectionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	res = &pb.DeleteCollectionRes{
		Code: http.StatusOK,
		Msg:  "delete collection successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
