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

	answerIndexModel := l.svcCtx.QuestionModel.AnswerIndex
	answerContentModel := l.svcCtx.QuestionModel.AnswerContent

	_, err = answerIndexModel.WithContext(l.ctx).
		Select(answerIndexModel.UserID, answerIndexModel.QuestionID).
		Where(answerIndexModel.UserID.Eq(j.Get("user_id").Int()),
			answerIndexModel.QuestionID.Eq(in.QuestionId),
		).
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

	answerIndex, err := answerIndexModel.WithContext(l.ctx).
		Where(answerIndexModel.QuestionID.Eq(in.QuestionId),
			answerIndexModel.UserID.Eq(j.Get("user_id").Int()),
			answerIndexModel.IPLoc.Eq(ip.GetIpLocFromApi(j.Get("last_ip").String()))).
		FirstOrCreate()
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
		Id:           answerIndex.ID,
		QuestionId:   answerIndex.QuestionID,
		UserId:       answerIndex.UserID,
		IpLoc:        answerIndex.IPLoc,
		ApproveCount: answerIndex.ApproveCount,
		LikeCount:    answerIndex.LikeCount,
		CollectCount: answerIndex.CollectCount,
		State:        answerIndex.State,
		Attrs:        answerIndex.Attrs,
		CreateTime:   answerIndex.CreateTime.String(),
		UpdateTime:   answerIndex.UpdateTime.String(),
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
		fmt.Sprintf("answerIndex_%d", answerIndex.ID),
		answerIndexBytes,
		time.Second*86400)

	// 发布消息-初始化评论模块
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
	err = commentMqProducer.DoCommentSubject(producer, 1, answerIndex.ID, "init")
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

	answerContent, err := answerContentModel.WithContext(l.ctx).
		Where(answerContentModel.AnswerID.Eq(answerIndex.ID),
			answerContentModel.Content.Eq(in.Content)).
		FirstOrCreate()
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
		AnswerId:   answerContent.AnswerID,
		Content:    answerContent.Content,
		Meta:       answerContent.Meta,
		CreateTime: answerIndex.CreateTime.String(),
		UpdateTime: answerIndex.UpdateTime.String(),
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
		fmt.Sprintf("answerContent_%d", answerContent.AnswerID),
		answerContentBytes,
		time.Second*86400)

	res = &pb.PublishAnswerRes{
		Code: http.StatusOK,
		Msg:  "publish answer successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
