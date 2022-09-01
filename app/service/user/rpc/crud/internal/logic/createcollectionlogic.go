package logic

import (
	"context"
	"main/app/common/log"
	"main/app/common/mq/nsq"
	notificationMqProducer "main/app/service/notification/mq/producer"
	questionMqProducer "main/app/service/question/mq/producer"
	"main/app/service/user/dao/model"
	userMqProducer "main/app/service/user/mq/producer"
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
			UserID:      in.UserId,
			CollectType: in.CollectType,
			ObjType:     in.ObjType,
			ObjID:       in.ObjId,
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

	// 发布通知
	producer, err := nsq.GetProducer()
	if err != nil {
		logger.Errorf("get producer failed, err: %v", err)
	} else {
		switch in.CollectType {
		case 1:
			// 喜欢
			err = notificationMqProducer.PublishNotification(producer, notificationMqProducer.PublishNotificationMessage{
				MessageType: 2,
				Data: notificationMqProducer.ApproveAndLikeData{
					Action:  1,
					ObjType: in.ObjType,
					ObjId:   in.ObjId,
				},
			})
			if err != nil {
				logger.Errorf("publish answer info to nsq failed, err: %v", err)
			}

		case 2:
			// 赞同
			err = notificationMqProducer.PublishNotification(producer, notificationMqProducer.PublishNotificationMessage{
				MessageType: 2,
				Data: notificationMqProducer.ApproveAndLikeData{
					UserId:  in.UserId,
					Action:  2,
					ObjType: in.ObjType,
					ObjId:   in.ObjId,
				},
			})
			if err != nil {
				logger.Errorf("publish answer info to nsq failed, err: %v", err)
			}

		case 3:
			// 收藏
			err = questionMqProducer.DoCollect(producer, questionMqProducer.CollectMessage{
				ObjType:  2,
				ObjId:    in.ObjId,
				AttrType: 2,
				Action:   0,
			})

		case 4:
			// 关注
			err = userMqProducer.ChangeFollower(producer, userMqProducer.ChangeFollowerMessage{
				UserId: in.ObjId,
				Action: 1,
			})
			if err != nil {
				logger.Errorf("publish user follower info to nsq failed,err: %v", err)
			}

			err = notificationMqProducer.PublishNotification(producer, notificationMqProducer.PublishNotificationMessage{
				MessageType: 1,
				Data:        notificationMqProducer.SubscriptionData{UserId: in.ObjId, FollowerId: in.UserId},
			})
			if err != nil {
				logger.Errorf("publish answer info to nsq failed, err: %v", err)
			}
		}
	}

	res = &pb.CreateCollectionRes{
		Code: http.StatusOK,
		Msg:  "create collection successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
