package logic

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"main/app/common/log"
	"main/app/common/mq/nsq"
	commentMqProducer "main/app/service/comment/mq/producer"
	notificationMqProducer "main/app/service/mq/nsq/producer/notification"
	"main/app/service/question/dao/model"
	modelpb "main/app/service/question/dao/pb"
	"main/app/service/question/rpc/crud/internal/svc"
	"main/app/service/question/rpc/crud/pb"
	"main/app/utils/net/ip"
	"net/http"
	"time"
)

type PublishAnswerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishAnswerLogic {
	return &PublishAnswerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishAnswerLogic) PublishAnswer(in *pb.PublishAnswerReq) (res *pb.PublishAnswerRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	j := gjson.Parse(in.UserDetails)

	userId := j.Get("user_id").Int()
	ipLoc := ip.GetIpLocFromApi(j.Get("last_ip").String())

	answerIndexModel := l.svcCtx.QuestionModel.AnswerIndex
	answerContentModel := l.svcCtx.QuestionModel.AnswerContent

	_, err = answerIndexModel.WithContext(l.ctx).
		Select(answerIndexModel.UserID, answerIndexModel.QuestionID).
		Where(answerIndexModel.UserID.Eq(userId),
			answerIndexModel.QuestionID.Eq(in.QuestionId)).
		First()
	switch err {
	case nil:
		res = &pb.PublishAnswerRes{
			Code: http.StatusForbidden,
			Msg:  "answer already exist",
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil

	case gorm.ErrRecordNotFound:

	default:
		logger.Errorf("publish answer failed, err: mysql err, %v", err)
		res = &pb.PublishAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	answerIndexId := l.svcCtx.IdGenerator.NewLong()
	nowTime := time.Now()

	err = answerIndexModel.WithContext(l.ctx).
		Create(&model.AnswerIndex{
			ID:         answerIndexId,
			QuestionID: in.QuestionId,
			UserID:     userId,
			IPLoc:      ipLoc,
			CreateTime: nowTime,
			UpdateTime: nowTime,
		})
	if err != nil {
		logger.Errorf("publish answer failed, err: %v", err)
		res = &pb.PublishAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	answerIndexProto := &modelpb.AnswerIndex{
		Id:         answerIndexId,
		QuestionId: in.QuestionId,
		UserId:     userId,
		IpLoc:      ipLoc,
		CreateTime: nowTime.String(),
		UpdateTime: nowTime.String(),
	}
	answerIndexBytes, err := proto.Marshal(answerIndexProto)
	if err != nil {
		logger.Errorf("marshal proto failed, err: %v", err)
		res = &pb.PublishAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
	l.svcCtx.Rdb.Set(l.ctx,
		fmt.Sprintf("answer_index_%d", answerIndexId),
		answerIndexBytes,
		time.Second*86400)

	// 发布消息-初始化评论模块 TODO: 使用asynq
	producer, err := nsq.GetProducer()
	if err != nil {
		logger.Errorf("get producer failed, err: %v", err)
		res = &pb.PublishAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", err)
		return res, nil
	}
	err = commentMqProducer.DoCommentSubject(producer, 1, answerIndexId, "init")
	if err != nil {
		logger.Errorf("publish answer info to nsq failed, err: %v", err)
		res = &pb.PublishAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", err)
		return res, nil
	}

	err = answerContentModel.WithContext(l.ctx).
		Create(&model.AnswerContent{
			AnswerID:   answerIndexId,
			Content:    in.Content,
			CreateTime: nowTime,
			UpdateTime: nowTime,
		})
	if err != nil {
		logger.Errorf("publish answer failed, err: mysql err, %v", err)
		res = &pb.PublishAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	answerContentProto := &modelpb.AnswerContent{
		AnswerId:   answerIndexId,
		Content:    in.Content,
		CreateTime: nowTime.String(),
		UpdateTime: nowTime.String(),
	}
	answerContentBytes, err := proto.Marshal(answerContentProto)
	if err != nil {
		logger.Errorf("marshal proto failed, err: %v", err)
		res = &pb.PublishAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
	l.svcCtx.Rdb.Set(l.ctx,
		fmt.Sprintf("answer_content_%d", answerIndexId),
		answerContentBytes,
		time.Second*86400)

	err = notificationMqProducer.PublishNotification(producer, notificationMqProducer.PublishNotificationMessage{
		MessageType: 5,
		Data: notificationMqProducer.AnswerData{
			UserId:     userId,
			QuestionId: in.QuestionId,
			AnswerId:   answerIndexId,
		},
	})
	if err != nil {
		logger.Errorf("publish msg to nsq failed, err: %v", err)
	}

	res = &pb.PublishAnswerRes{
		Code: http.StatusOK,
		Msg:  "publish answer successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
