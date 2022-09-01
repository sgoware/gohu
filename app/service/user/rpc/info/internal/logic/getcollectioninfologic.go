package logic

import (
	"context"
	"main/app/common/log"
	"net/http"

	"main/app/service/user/rpc/info/internal/svc"
	"main/app/service/user/rpc/info/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCollectionInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCollectionInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCollectionInfoLogic {
	return &GetCollectionInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCollectionInfoLogic) GetCollectionInfo(in *pb.GetCollectionInfoReq) (res *pb.GetCollectionInfoRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	userCollectModel := l.svcCtx.UserModel.UserCollection

	userCollections, err := userCollectModel.WithContext(l.ctx).
		Where(userCollectModel.UserID.Eq(in.UserId), userCollectModel.CollectType.Eq(in.CollectionType)).Find()
	if err != nil {
		logger.Debugf("get collections failed, err: mysql err, %v", err)
		res = &pb.GetCollectionInfoRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
			Data: nil,
		}
		return res, nil
	}

	res = &pb.GetCollectionInfoRes{
		Code: http.StatusOK,
		Msg:  "get collections successfully",
		Ok:   true,
		Data: &pb.GetCollectionInfoRes_Data{
			ObjType: make([]int32, 0),
			ObjId:   make([]int64, 0),
		},
	}
	for _, userCollection := range userCollections {
		res.Data.ObjType = append(res.Data.ObjType, userCollection.ObjType)
		res.Data.ObjId = append(res.Data.ObjId, userCollection.ObjID)
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
