package logic

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
	"main/app/common/log"
	"main/app/service/mq/asynq/processor/job"
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

	payload, err := json.Marshal(&job.MsgCrudCommentSubjectPayload{
		Action: 3,
		Id:     in.AnswerId,
	})
	if err != nil {
		logger.Errorf("marshal [MsgCrudCommentSubjectPayload] failed, err: %v", err)
		res = &pb.DeleteAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", err)
		return res, nil
	}

	_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(job.MsgCrudCommentSubjectTask, payload))
	if err != nil {
		logger.Errorf("create [MsgCrudCommentSubjectTask] insert queue failed, err: %v", err)
		res = &pb.DeleteAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", err)
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
