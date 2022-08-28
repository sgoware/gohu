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

type UpdateAnswerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAnswerLogic {
	return &UpdateAnswerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAnswerLogic) UpdateAnswer(in *pb.UpdateAnswerReq) (res *pb.UpdateAnswerRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	answerContentModel := l.svcCtx.QuestionModel.AnswerContent

	_, err = answerContentModel.WithContext(l.ctx).Select(answerContentModel.Content).
		Where(answerContentModel.AnswerID.Eq(in.AnswerId)).
		Update(answerContentModel.Content, in.Content)
	switch err {
	case gorm.ErrRecordNotFound:
		res = &pb.UpdateAnswerRes{
			Code: http.StatusBadRequest,
			Msg:  "answer not found",
			Ok:   false,
		}
	case nil:
		res = &pb.UpdateAnswerRes{
			Code: http.StatusOK,
			Msg:  "update answer successfully",
			Ok:   true,
		}
	default:
		logger.Errorf("update question failed, err: %v", err)
		res = &pb.UpdateAnswerRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
	}

	logger.Debugf("send message: %v", res.String())
	return res, nil
}
