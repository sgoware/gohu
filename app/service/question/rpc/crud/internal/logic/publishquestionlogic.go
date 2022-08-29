package logic

import (
	"context"
	"main/app/common/log"
	"main/app/service/question/dao/model"
	"main/app/service/question/rpc/crud/internal/svc"
	"main/app/service/question/rpc/crud/pb"
	"main/app/utils/net/ip"
	"net/http"

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

	questionSubjectModel.WithContext(l.ctx)
	err = questionSubjectModel.WithContext(l.ctx).Create(&model.QuestionSubject{
		UserID: j.Get("user_id").Int(),
		IPLoc:  ip.GetIpLocFromApi(j.Get("last_ip").String()),
		Title:  in.Title,
		Topic:  in.Topic,
		Tag:    in.Tag,
	})
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

	questionSubject, err := questionSubjectModel.WithContext(l.ctx).
		Select(questionSubjectModel.ID, questionSubjectModel.UserID).
		Where(questionSubjectModel.UserID.Eq(j.Get("user_id").Int())).
		Order(questionSubjectModel.ID.Desc()).Last()
	if err != nil {
		logger.Errorf("publish question failed, err: mysql err, %v", err)
		res = &pb.PublishQuestionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	err = questionContentModel.WithContext(l.ctx).Create(&model.QuestionContent{
		QuestionID: questionSubject.ID,
		Content:    in.Content,
	})
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

	res = &pb.PublishQuestionRes{
		Code: http.StatusOK,
		Msg:  "publish question successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
