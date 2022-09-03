package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"google.golang.org/protobuf/proto"
	"main/app/common/log"
	"main/app/common/mq/nsq"
	"main/app/service/mq/asynq/processor/job"
	notificationMqProducer "main/app/service/mq/nsq/producer/notification"
	modelpb "main/app/service/question/dao/pb"
	"main/app/service/question/rpc/crud/internal/svc"
	"main/app/service/question/rpc/crud/pb"
	"main/app/utils/net/ip"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
)

type PublishQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishQuestionLogic {
	return &PublishQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishQuestionLogic) PublishQuestion(in *pb.PublishQuestionReq) (res *pb.PublishQuestionRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	j := gjson.Parse(in.UserDetails)

	userId := j.Get("user_id").Int()

	questionSubjectId := l.svcCtx.IdGenerator.NewLong()

	nowTime := time.Now()

	ipLoc := ip.GetIpLocFromApi(j.Get("last_ip").String())

	questionSubjectProto := &modelpb.QuestionSubject{
		Id:         questionSubjectId,
		UserId:     userId,
		IpLoc:      ipLoc,
		Title:      in.Title,
		Topic:      in.Topic,
		Tag:        in.Tag,
		CreateTime: nowTime.String(),
		UpdateTime: nowTime.String(),
	}
	bytes, err := proto.Marshal(questionSubjectProto)
	if err != nil {
		logger.Errorf("marshal proto failed, err: %v", err)
		res = &pb.PublishQuestionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
	l.svcCtx.Rdb.Set(l.ctx,
		fmt.Sprintf("question_subject_%d", questionSubjectId),
		bytes,
		time.Second*86400)

	payload, err := json.Marshal(job.MsgCrudQuestionSubjectRecordPayload{
		Action:     1,
		Id:         questionSubjectId,
		UserId:     userId,
		IpLoc:      ipLoc,
		Title:      in.Title,
		Topic:      in.Topic,
		Tag:        in.Tag,
		CreateTime: nowTime,
		UpdateTime: nowTime,
	})
	if err != nil {
		logger.Errorf("marshal [MsgCrudQuestionSubjectRecordPayload] failed, err: %v", err)
		res = &pb.PublishQuestionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", err)
		return res, nil
	}

	_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(job.MsgCrudQuestionSubjectRecordTask, payload))
	if err != nil {
		logger.Errorf("create [MsgCrudQuestionSubjectRecordTask] insert queue failed, err: %v", err)
		res = &pb.PublishQuestionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", err)
		return res, nil
	}

	questionContentProto := &modelpb.QuestionContent{
		QuestionId: questionSubjectId,
		Content:    in.Content,
		Meta:       "",
		CreateTime: nowTime.String(),
		UpdateTime: nowTime.String(),
	}
	bytes, err = proto.Marshal(questionContentProto)
	if err != nil {
		logger.Errorf("marshal proto failed, err: %v", err)
		res = &pb.PublishQuestionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}
	l.svcCtx.Rdb.Set(l.ctx,
		fmt.Sprintf("question_content_%d", questionSubjectId),
		bytes,
		time.Second*86400)

	payload, err = json.Marshal(job.MsgCrudQuestionContentRecordPayload{
		Action:     1,
		QuestionId: questionSubjectId,
		Content:    in.Content,
		CreateTime: nowTime,
		UpdateTime: nowTime,
	})
	if err != nil {
		logger.Errorf("marshal [MsgCrudQuestionContentRecordPayload] failed, err: %v", err)
		res = &pb.PublishQuestionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", err)
		return res, nil
	}

	_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(job.MsgCrudQuestionContentRecordTask, payload))
	if err != nil {
		logger.Errorf("create [MsgCrudQuestionContentRecordTask] insert queue failed, err: %v", err)
		res = &pb.PublishQuestionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", err)
		return res, nil
	}

	producer, err := nsq.GetProducer()
	if err != nil {
		logger.Errorf("get nsq producer failed, err: %v", err)
	} else {
		err = notificationMqProducer.PublishNotification(producer, notificationMqProducer.PublishNotificationMessage{
			MessageType: 4,
			Data: notificationMqProducer.SubscriptionData{
				UserId:  userId,
				Action:  1,
				ObjType: 1,
				ObjId:   questionSubjectId,
			},
		})
		if err != nil {
			logger.Errorf("publish msg to nsq failed, err: %v", err)
		}
	}

	err = l.svcCtx.Rdb.SAdd(l.ctx,
		fmt.Sprintf("question_id_user_set_%d", userId),
		questionSubjectId).Err()
	if err != nil {
		logger.Errorf("update [question_id_user_set] failed, err: %v", err)
	}

	res = &pb.PublishQuestionRes{
		Code: http.StatusOK,
		Msg:  "publish question successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
