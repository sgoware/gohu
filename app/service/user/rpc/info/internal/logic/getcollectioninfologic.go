package logic

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"main/app/common/log"
	"net/http"
	"strings"

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

	userCollectionsCache, err := l.svcCtx.Rdb.SMembers(l.ctx,
		fmt.Sprintf("user_collect_%d_%d", in.UserId, in.CollectionType)).Result()
	if err == nil {
		objTypes := make([]int32, 0)
		objIds := make([]int64, 0)
		for _, userCollectionCache := range userCollectionsCache {
			output := strings.Split(userCollectionCache, ":")
			objTypes = append(objTypes, cast.ToInt32(output[0]))
			objIds = append(objIds, cast.ToInt64(output[1]))
		}
		res = &pb.GetCollectionInfoRes{
			Code: http.StatusOK,
			Msg:  "get collections successfully",
			Ok:   true,
			Data: &pb.GetCollectionInfoRes_Data{
				ObjType: objTypes,
				ObjId:   objIds,
			},
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
	logger.Errorf("get [user_collect] cache failed, err: %v", err)

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
		l.svcCtx.Rdb.SAdd(l.ctx,
			fmt.Sprintf("user_collect_%d_%d", in.UserId, in.CollectionType),
			fmt.Sprintf("%d:%d", userCollection.ObjType, userCollection.ObjID))
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
