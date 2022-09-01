package logic

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"main/app/common/log"
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

	questionSubjectModel := l.svcCtx.QuestionModel.QuestionSubject
	questionContentModel := l.svcCtx.QuestionModel.QuestionContent

	questionSubject, err := questionSubjectModel.WithContext(l.ctx).
		Where(questionSubjectModel.UserID.Eq(j.Get("user_id").Int()),
			questionSubjectModel.IPLoc.Eq(ip.GetIpLocFromApi(j.Get("last_ip").String())),
			questionSubjectModel.Title.Eq(in.Title),
			questionSubjectModel.Topic.Eq(in.Topic),
			questionSubjectModel.Tag.Eq(in.Tag)).
		FirstOrCreate()
	if err != nil {
		logger.Errorf("publish question failed, err: %v", err)
		res = &pb.PublishQuestionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	questionSubjectProto := &modelpb.QuestionSubject{
		Id:          questionSubject.ID,
		UserId:      questionSubject.UserID,
		IpLoc:       questionSubject.IPLoc,
		Title:       questionSubject.Title,
		Topic:       questionSubject.Topic,
		Tag:         questionSubject.Tag,
		SubCount:    questionSubject.SubCount,
		AnswerCount: questionSubject.AnswerCount,
		ViewCount:   questionSubject.ViewCount,
		State:       questionSubject.State,
		CreateTime:  questionSubject.CreateTime.String(),
		UpdateTime:  questionSubject.UpdateTime.String(),
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
		fmt.Sprintf("questionSubject_%d", questionSubject.ID),
		bytes,
		time.Second*86400)

	questionContent, err := questionContentModel.WithContext(l.ctx).
		Where(questionContentModel.QuestionID.Eq(questionSubject.ID),
			questionContentModel.Content.Eq(in.Content)).
		FirstOrCreate()
	if err != nil {
		logger.Errorf("publish question failed, err: %v", err)
		res = &pb.PublishQuestionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	questionContentProto := &modelpb.QuestionContent{
		QuestionId: questionContent.QuestionID,
		Content:    questionContent.Content,
		Meta:       questionContent.Meta,
		CreateTime: questionContent.CreateTime.String(),
		UpdateTime: questionContent.UpdateTime.String(),
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
		fmt.Sprintf("questionContent_%d", questionContent.QuestionID),
		bytes,
		time.Second*86400)

	res = &pb.PublishQuestionRes{
		Code: http.StatusOK,
		Msg:  "publish question successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
