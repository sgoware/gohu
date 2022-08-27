package logic

import (
	"context"
	"gorm.io/gorm"
	"main/app/common/log"
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
			Mag:  "delete answer failed, err: answer not found",
			Ok:   false,
		}
	case nil:

	default:
		logger.Errorf("update Answer failed, err: %v", err)
		res = &pb.DeleteAnswerRes{
			Code: http.StatusInternalServerError,
			Mag:  "delete answer failed, err: internal err",
			Ok:   false,
		}
	}

	res = &pb.DeleteAnswerRes{
		Code: http.StatusOK,
		Mag:  "delete answer successfully",
		Ok:   true,
	}

	logger.Debugf("send message: %v", res.String())
	return res, nil
}
