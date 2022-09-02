package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"main/app/common/log"
	"main/app/common/mq/nsq"
	"main/app/service/mq/asynq/processor/job"
	notificationMqProducer "main/app/service/mq/nsq/producer/notification"
	questionMqProducer "main/app/service/question/mq/producer"
	"main/app/service/user/rpc/crud/internal/svc"
	"main/app/service/user/rpc/crud/pb"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type DoCollectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDoCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DoCollectionLogic {
	return &DoCollectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DoCollectionLogic) DoCollection(in *pb.DoCollectionReq) (res *pb.DoCollectionRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	userCollectionModel := l.svcCtx.UserModel.UserCollection

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
			ok, err := l.svcCtx.Rdb.SIsMember(l.ctx,
				fmt.Sprintf("user_collect_%d_%d_%d", in.UserId, 4, in.ObjType),
				in.ObjId).Result()
			if err == nil {
				if !ok {
					// 不存在收藏的缓存, 则是创建操作
					err = createCollectionCache(l.ctx, l.svcCtx, in)

					switch in.ObjType {
					case 1:
						// 关注用户
						// 通知被关注的用户
						err = notificationMqProducer.PublishNotification(producer, notificationMqProducer.PublishNotificationMessage{
							MessageType: 1,
							Data:        notificationMqProducer.SubscriptionData{UserId: in.ObjId, FollowerId: in.UserId},
						})
						if err != nil {
							logger.Errorf("publish notificaion to nsq failed, %v", err)
							res = &pb.DoCollectionRes{
								Code: http.StatusInternalServerError,
								Msg:  "internal err",
								Ok:   false,
							}
							logger.Debugf("send message: %v", err)
							return res, nil
						}
						// 更新 user_subject 缓存
						payload, err := json.Marshal(&job.MsgAddUserSubjectCachePayload{Id: in.ObjId, Follower: 1})
						if err != nil {
							logger.Errorf("marshal [MsgAddUserSubjectCachePayload] failed, %v", err)
							res = &pb.DoCollectionRes{
								Code: http.StatusInternalServerError,
								Msg:  "internal err",
								Ok:   false,
							}
							logger.Debugf("send message: %v", err)
							return res, nil
						}

						_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(job.MsgAddUserSubjectCacheTask, payload))
						if err != nil {
							logger.Errorf("create [MsgAddUserSubjectCacheTask] insert queue failed, %v", err)
							res = &pb.DoCollectionRes{
								Code: http.StatusInternalServerError,
								Msg:  "internal err",
								Ok:   false,
							}
							logger.Debugf("send message: %v", err)
							return res, nil
						}

						// 关注者计数器+1, 队列调度器定时更新数据库
						err = l.svcCtx.Rdb.SAdd(l.ctx,
							"user_follower",
							in.ObjId).Err()
						if err != nil {
							logger.Errorf("add [user_follower] member failed, err: %v", err)
							res = &pb.DoCollectionRes{
								Code: http.StatusInternalServerError,
								Msg:  "internal err",
								Ok:   false,
							}
							logger.Debugf("send message: %v", err)
							return res, nil
						}
						err = l.svcCtx.Rdb.Incr(l.ctx,
							fmt.Sprintf("user_follower_%d", in.ObjId)).Err()
						if err != nil {
							logger.Errorf("increase [user_follower] failed, err: %v", err)
							res = &pb.DoCollectionRes{
								Code: http.StatusInternalServerError,
								Msg:  "internal err",
								Ok:   false,
							}
							logger.Debugf("send message: %v", err)
							return res, nil
						}
						res = &pb.DoCollectionRes{
							Code: http.StatusOK,
							Msg:  "follow user successfully",
							Ok:   true,
						}
						logger.Debugf("send message: %v", err)
						return res, nil
					}
				} else {
					// 存在收藏的缓存, 则是删除操作
					err = deleteCollectionCache(l.ctx, l.svcCtx, in)

					switch in.ObjType {
					case 1:
						// 取消关注用户
						// 更新 user_subject 缓存
						payload, err := json.Marshal(&job.MsgAddUserSubjectCachePayload{Id: in.ObjId, Follower: -1})
						if err != nil {
							logger.Errorf("marshal [MsgAddUserSubjectCachePayload] failed, %v", err)
							res = &pb.DoCollectionRes{
								Code: http.StatusInternalServerError,
								Msg:  "internal err",
								Ok:   false,
							}
							logger.Debugf("send message: %v", err)
							return res, nil
						}

						_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(job.MsgAddUserSubjectCacheTask, payload))
						if err != nil {
							logger.Errorf("create [MsgAddUserSubjectCacheTask] insert queue failed, %v", err)
							res = &pb.DoCollectionRes{
								Code: http.StatusInternalServerError,
								Msg:  "internal err",
								Ok:   false,
							}
							logger.Debugf("send message: %v", err)
							return res, nil
						}

						// 关注者计数器-1, 队列调度器定时更新数据库
						err = l.svcCtx.Rdb.SRem(l.ctx,
							"user_follower",
							in.ObjId).Err()
						if err != nil {
							logger.Errorf("add [user_follower] member failed, err: %v", err)
							res = &pb.DoCollectionRes{
								Code: http.StatusInternalServerError,
								Msg:  "internal err",
								Ok:   false,
							}
							logger.Debugf("send message: %v", err)
							return res, nil
						}
						err = l.svcCtx.Rdb.Decr(l.ctx,
							fmt.Sprintf("user_follower_%d", in.UserId)).Err()
						if err != nil {
							logger.Errorf("increase [user_follower] failed, err: %v", err)
							res = &pb.DoCollectionRes{
								Code: http.StatusInternalServerError,
								Msg:  "internal err",
								Ok:   false,
							}
							logger.Debugf("send message: %v", err)
							return res, nil
						}
						if err != nil {
							logger.Errorf("unfollow user failed, err: %v", err)
							res = &pb.DoCollectionRes{
								Code: http.StatusInternalServerError,
								Msg:  "internal err",
								Ok:   false,
							}
							logger.Debugf("send message: %v", err)
							return res, nil
						}
						res = &pb.DoCollectionRes{
							Code: http.StatusOK,
							Msg:  "unfollow user successfully",
							Ok:   true,
						}
						logger.Debugf("send message: %v", err)
						return res, nil
					}
				}
			} else {
				logger.Errorf("get [user_collect] cache member failed, err: %v", err)
			}
			// 缓存获取失败, 看看数据库
			cnt, err := userCollectionModel.WithContext(l.ctx).
				Select(userCollectionModel.UserID.Eq(in.UserId),
					userCollectionModel.ObjType.Eq(in.ObjType),
					userCollectionModel.ObjID.Eq(in.ObjId)).
				Count()
			if err != nil {
				logger.Errorf("query [user_collect] record failed, err: %v", err)
				res = &pb.DoCollectionRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				logger.Debugf("send message: %v", err)
				return res, nil
			}
			if cnt == 0 {
				// 不存在, 则是创建操作
				err = createCollectionCache(l.ctx, l.svcCtx, in)
				if err != nil {
					logger.Errorf("follow user failed, err: %v", err)
					res = &pb.DoCollectionRes{
						Code: http.StatusInternalServerError,
						Msg:  "internal err",
						Ok:   false,
					}
					return res, nil
				}
				switch in.ObjType {
				case 1:
					// 关注用户

					// 通知用户被关注了
					err = notificationMqProducer.PublishNotification(producer, notificationMqProducer.PublishNotificationMessage{
						MessageType: 1,
						Data:        notificationMqProducer.SubscriptionData{UserId: in.ObjId, FollowerId: in.UserId},
					})
					if err != nil {
						logger.Errorf("publish notification info to nsq failed, %v", err)
						res = &pb.DoCollectionRes{
							Code: http.StatusInternalServerError,
							Msg:  "internal err",
							Ok:   false,
						}
						logger.Debugf("send message: %v", err)
						return res, nil
					}

					res = &pb.DoCollectionRes{
						Code: http.StatusOK,
						Msg:  "follow user successfully",
						Ok:   true,
					}
					logger.Debugf("send message: %v", err)
					return res, nil
				}
			}
			// 存在记录, 则是删除操作
			err = deleteCollectionCache(l.ctx, l.svcCtx, in)
			if err != nil {
				logger.Errorf("unfollow user failed, err: %v", err)
				res = &pb.DoCollectionRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				return res, nil
			}
			res = &pb.DoCollectionRes{
				Code: http.StatusOK,
				Msg:  "unfollow user successfully",
				Ok:   true,
			}
			logger.Debugf("send message: %v", err)
			return res, nil
		}
	}
	return nil, nil
}

func createCollectionCache(ctx context.Context, svcCtx *svc.ServiceContext, in *pb.DoCollectionReq) (err error) {
	// 更新 user_collect 缓存
	err = svcCtx.Rdb.SAdd(ctx,
		fmt.Sprintf("user_collect_%d_%d_%d", in.UserId, in.CollectType, in.ObjType),
		in.ObjId).Err()
	if err != nil {
		return fmt.Errorf("add [user_collect] cache member failed, %v", err)
	}

	err = svcCtx.Rdb.LPush(ctx,
		"user_collect",
		fmt.Sprintf("0_%d_%d_%d_%d", in.UserId, in.CollectType, in.ObjType, in.ObjId)).Err()
	if err != nil {
		return fmt.Errorf("push [user_collect] cache failed, %v", err)
	}

	return nil
}

func deleteCollectionCache(ctx context.Context, svcCtx *svc.ServiceContext, in *pb.DoCollectionReq) (err error) {
	// 更新 user_collect 缓存
	err = svcCtx.Rdb.SRem(ctx,
		fmt.Sprintf("user_collect_%d_%d_%d", in.UserId, in.CollectType, in.ObjType),
	).Err()
	if err != nil {
		return fmt.Errorf("delete [user_collect] cache member failed, %v", err)
	}

	err = svcCtx.Rdb.LPush(ctx,
		"user_collect",
		fmt.Sprintf("1_%d_%d_%d_%d", in.UserId, in.CollectType, in.ObjType, in.ObjId)).Err()
	if err != nil {
		return fmt.Errorf("push [user_collect] cache failed, %v", err)
	}

	return nil
}
