package logic

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"main/app/common/log"
	"main/app/service/user/rpc/info/internal/svc"
	"main/app/service/user/rpc/info/pb"
	"net/http"

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
		fmt.Sprintf("user_collect_set_%d_%d_%d", in.UserId, in.CollectionType, in.ObjType)).Result()
	if err == nil {
		if len(userCollectionsCache) > 1 {
			objIds := make([]int64, 0)
			for _, userCollectionCache := range userCollectionsCache {
				if userCollectionCache != "0" {
					objIds = append(objIds, cast.ToInt64(userCollectionCache))
				}
			}
			res = &pb.GetCollectionInfoRes{
				Code: http.StatusOK,
				Msg:  "get collections successfully",
				Ok:   true,
				Data: &pb.GetCollectionInfoRes_Data{
					ObjId: objIds,
				},
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
	}
	logger.Errorf("get [user_collect] cache failed, err: %v", err)

	l.svcCtx.Rdb.SAdd(l.ctx,
		fmt.Sprintf("user_collect_set_%d_%d_%d", in.UserId, in.CollectionType, in.ObjType),
		0)

	userCollectModel := l.svcCtx.UserModel.UserCollection

	userCollections, err := userCollectModel.WithContext(l.ctx).
		Where(userCollectModel.UserID.Eq(in.UserId),
			userCollectModel.CollectType.Eq(in.CollectionType),
			userCollectModel.ObjType.Eq(in.ObjType)).Find()
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
			ObjId: make([]int64, 0),
		},
	}
	for _, userCollection := range userCollections {
		res.Data.ObjId = append(res.Data.ObjId, userCollection.ObjID)
		l.svcCtx.Rdb.SAdd(l.ctx,
			fmt.Sprintf("user_collect_set_%d_%d_%d", in.UserId, in.CollectionType, in.ObjType),
			userCollection.ObjID)
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
