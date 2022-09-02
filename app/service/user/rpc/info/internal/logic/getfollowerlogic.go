package logic

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"main/app/common/log"
	"net/http"

	"main/app/service/user/rpc/info/internal/svc"
	"main/app/service/user/rpc/info/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowerLogic {
	return &GetFollowerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowerLogic) GetFollower(in *pb.GetFollowerReq) (res *pb.GetFollowerRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	followerIds := make([]int64, 0)

	followerMembersCache, err := l.svcCtx.Rdb.SMembers(l.ctx,
		fmt.Sprintf("user_follower_member_%d", in.UserId)).Result()
	if err == nil {
		if len(followerMembersCache) > 1 {
			for _, followerMemberCache := range followerMembersCache {
				followerIds = append(followerIds, cast.ToInt64(followerMemberCache))
			}
			res = &pb.GetFollowerRes{
				Code: http.StatusOK,
				Msg:  "get follower ids successfully",
				Ok:   true,
				Data: &pb.GetFollowerRes_Data{UserIds: followerIds},
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
	} else {
		logger.Errorf("get [user_follower_member] cacche failed, err: %v", err)
	}

	l.svcCtx.Rdb.SAdd(l.ctx,
		fmt.Sprintf("user_follower_member_%d", in.UserId),
		0)

	userCollectionModel := l.svcCtx.UserModel.UserCollection

	userCollections, err := userCollectionModel.WithContext(l.ctx).
		Select(userCollectionModel.UserID,
			userCollectionModel.CollectType,
			userCollectionModel.ObjType,
			userCollectionModel.ObjID).
		Where(userCollectionModel.CollectType.Eq(4),
			userCollectionModel.ObjType.Eq(1),
			userCollectionModel.ObjID.Eq(in.UserId)).
		Find()
	if err != nil {
		logger.Errorf("query [user_collection] record failed, err: %v", err)
		res = &pb.GetFollowerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	for _, userCollection := range userCollections {
		followerIds = append(followerIds, userCollection.UserID)
		l.svcCtx.Rdb.SAdd(l.ctx,
			fmt.Sprintf("user_follower_member_%d", in.UserId),
			userCollection.UserID)
	}

	res = &pb.GetFollowerRes{
		Code: http.StatusOK,
		Msg:  "get follower ids successfully",
		Ok:   true,
		Data: &pb.GetFollowerRes_Data{UserIds: followerIds},
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
