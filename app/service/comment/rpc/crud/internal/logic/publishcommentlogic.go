package logic

import (
	"context"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"main/app/common/log"
	"main/app/common/mq/nsq"
	"main/app/service/comment/dao/model"
	"main/app/service/comment/rpc/crud/internal/svc"
	"main/app/service/comment/rpc/crud/pb"
	notificationMqProducer "main/app/service/mq/nsq/producer/notification"
	"main/app/utils/net/ip"
	"net/http"
	"time"
)

type PublishCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishCommentLogic {
	return &PublishCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishCommentLogic) PublishComment(in *pb.PublishCommentReq) (res *pb.PublishCommentRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	j := gjson.Parse(in.UserDetails)
	userId := j.Get("user_id").Int()

	commentId := l.svcCtx.IdGenerator.NewLong()
	nowTime := time.Now()
	ipLoc := ip.GetIpLocFromApi(j.Get("last_ip").String())

	commentIndexModel := l.svcCtx.CommentModel.CommentIndex
	commentContentModel := l.svcCtx.CommentModel.CommentContent

	if in.RootId == 0 {
		// 是评论的情况
		count, err := commentIndexModel.WithContext(l.ctx).
			Where(commentIndexModel.SubjectID.Eq(in.SubjectId),
				commentIndexModel.RootID.Eq(0)).
			Count()
		if err != nil {
			logger.Errorf("publish comment failed, err: mysql err, %v", err)
			res = &pb.PublishCommentRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		err = commentIndexModel.WithContext(l.ctx).
			Create(&model.CommentIndex{
				ID:           commentId,
				SubjectID:    in.SubjectId,
				UserID:       userId,
				IPLoc:        ipLoc,
				CommentFloor: int32(count + 1),
				CommentID:    0,
				ReplyFloor:   0,
				ApproveCount: 0,
				State:        0,
				Attrs:        0,
				CreateTime:   nowTime,
				UpdateTime:   nowTime,
			})
		if err != nil {
			logger.Errorf("publish comment failed, err: mysql err, %v", err)
			res = &pb.PublishCommentRes{
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
			err = notificationMqProducer.PublishNotification(producer, notificationMqProducer.PublishNotificationMessage{
				MessageType: 3,
				Data: notificationMqProducer.CommentData{
					UserId:    userId,
					SubjectId: in.SubjectId,
					CommentId: 0,
				},
			})
		}
	} else {
		// 是回复评论的情况
		count, err := commentIndexModel.WithContext(l.ctx).
			Where(commentIndexModel.SubjectID.Eq(in.SubjectId),
				commentIndexModel.RootID.Eq(in.RootId)).
			Count()
		if err != nil {
			logger.Errorf("publish comment failed, err: mysql err, %v", err)
			res = &pb.PublishCommentRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		err = commentIndexModel.WithContext(l.ctx).
			Create(&model.CommentIndex{
				ID:         commentId,
				SubjectID:  in.SubjectId,
				UserID:     userId,
				IPLoc:      ipLoc,
				RootID:     in.RootId,
				CommentID:  in.CommentId,
				ReplyFloor: int32(count + 1),
				CreateTime: nowTime,
				UpdateTime: nowTime,
			})
		if err != nil {
			logger.Errorf("publish comment failed, err: mysql err, %v", err)
			res = &pb.PublishCommentRes{
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
			err = notificationMqProducer.PublishNotification(producer, notificationMqProducer.PublishNotificationMessage{
				MessageType: 3,
				Data: notificationMqProducer.CommentData{
					UserId:    userId,
					SubjectId: in.SubjectId,
					CommentId: in.CommentId,
				},
			})
		}
	}

	err = commentContentModel.WithContext(l.ctx).
		Create(&model.CommentContent{
			CommentID: commentId,
			Content:   in.Content,
		})
	if err != nil {
		logger.Errorf("publish comment failed, err: mysql err, %v", err)
		res = &pb.PublishCommentRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	res = &pb.PublishCommentRes{
		Code: http.StatusOK,
		Msg:  "publish comment successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
