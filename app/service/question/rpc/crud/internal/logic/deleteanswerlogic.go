package logic

import (
	"context"
	"gorm.io/gorm"
	"main/app/common/log"
	"main/app/common/mq/nsq"
	commentMqProducer "main/app/service/comment/mq/producer"
	"net/http"

	"main/app/service/question/rpc/crud/internal/svc"
	"main/app/service/question/rpc/crud/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAnswerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAnswerLogic {
	return &DeleteAnswerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAnswerLogic) DeleteAnswer(in *pb.DeleteAnswerReq) (res *pb.DeleteAnswerRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	// 删除answer_index后, 级联删除关联的回答内容和评论
	answerIndexModel := l.svcCtx.QuestionModel.AnswerIndex

	_, err = answerIndexModel.WithContext(l.ctx).Where(answerIndexModel.ID.Eq(in.AnswerId)).Delete()
	switch err {
	case gorm.ErrRecordNotFound:
		res = &pb.DeleteAnswerRes{
			Code: http.StatusBadRequest,
			Msg:  "answer not found",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	case nil:

	default:
		logger.Errorf("update Answer failed, err: %v", err)
		res = &pb.DeleteAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	// 发布消息-删除评论模块
	producer, err := nsq.GetProducer()
	err = commentMqProducer.DoCommentSubject(producer, 1, in.AnswerId, "delete")
	if err != nil {
		logger.Errorf("publish answer info to nsq failed, err: %v", err)
		res = &pb.DeleteAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		return res, nil
	}

	res = &pb.DeleteAnswerRes{
		Code: http.StatusOK,
		Msg:  "delete answer successfully",
		Ok:   true,
	}

	logger.Debugf("send message: %v", res.String())
	return res, nil
}
