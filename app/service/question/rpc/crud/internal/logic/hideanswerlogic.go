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

type HideAnswerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHideAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HideAnswerLogic {
	return &HideAnswerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HideAnswerLogic) HideAnswer(in *pb.HideAnswerReq) (res *pb.HideAnswerRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	answerIndexModel := l.svcCtx.QuestionModel.AnswerIndex

	_, err = answerIndexModel.WithContext(l.ctx).Select(answerIndexModel.State).
		Where(answerIndexModel.ID.Eq(in.AnswerId)).
		Update(answerIndexModel.State, 1)
	switch err {
	case gorm.ErrRecordNotFound:
		res = &pb.HideAnswerRes{
			Code: http.StatusBadRequest,
			Mag:  "hide answer failed, err: answer not found",
			Ok:   false,
		}
	case nil:
		res = &pb.HideAnswerRes{
			Code: http.StatusOK,
			Mag:  "hide answer successfully",
			Ok:   true,
		}
	default:
		logger.Errorf("update question failed, err: %v", err)
		res = &pb.HideAnswerRes{
			Code: http.StatusInternalServerError,
			Mag:  "hide answer failed, err: internal err",
			Ok:   false,
		}
	}

	logger.Debugf("send message: %v", res.String())
	return res, nil
}
