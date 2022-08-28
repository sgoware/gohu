package logic

import (
	"context"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"main/app/common/log"
	"main/app/service/question/dao/model"
	"main/app/service/question/rpc/crud/internal/svc"
	"main/app/service/question/rpc/crud/pb"
	"main/app/utils/net/ip"
	"net/http"
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

	err = answerIndexModel.WithContext(l.ctx).Create(&model.AnswerIndex{
		QuestionID: in.QuestionId,
		UserID:     j.Get("user_id").Int(),
	})
	if err != nil {
		logger.Errorf("publish answer failed, err: %v", err)
		res = &pb.PublishAnswerRes{
			Code: http.StatusInternalServerError,
			Mag:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return nil, err
	}

	answerIndex, err := answerIndexModel.WithContext(l.ctx).
		Select(answerIndexModel.ID, answerIndexModel.UserID).
		Where(answerIndexModel.UserID.Eq(j.Get("user_id").Int())).
		Order(answerIndexModel.UserID.Desc()).Last()
	if err != nil {
		logger.Errorf("publish answer failed, err: mysql err, %v", err)
		res = &pb.PublishAnswerRes{
			Code: http.StatusInternalServerError,
			Mag:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	err = answerContentModel.WithContext(l.ctx).Create(&model.AnswerContent{
		AnswerID: answerIndex.ID,
		Content:  in.Content,
		IPLoc:    ip.GetIpLocFromApi(j.Get("last_ip").String()),
	})
	if err != nil {
		logger.Errorf("publish answer failed, err: %v", err)
		res = &pb.PublishAnswerRes{
			Code: http.StatusInternalServerError,
			Mag:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return nil, err
	}

	logger.Debugf("send message: %v", res.String())
	return res, nil
}
