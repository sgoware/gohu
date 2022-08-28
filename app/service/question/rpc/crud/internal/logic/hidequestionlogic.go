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

type HideQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHideQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HideQuestionLogic {
	return &HideQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HideQuestionLogic) HideQuestion(in *pb.HideQuestionReq) (res *pb.HideQuestionRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	questionSubjectModel := l.svcCtx.QuestionModel.QuestionSubject

	_, err = questionSubjectModel.WithContext(l.ctx).Select(questionSubjectModel.State).
		Where(questionSubjectModel.ID.Eq(in.QuestionId)).
		Update(questionSubjectModel.State, 1)
	switch err {
	case gorm.ErrRecordNotFound:
		res = &pb.HideQuestionRes{
			Code: http.StatusBadRequest,
			Mag:  "question not found",
			Ok:   false,
		}
	case nil:
		res = &pb.HideQuestionRes{
			Code: http.StatusOK,
			Mag:  "hide question successfully",
			Ok:   true,
		}
	default:
		logger.Errorf("update question failed, err: %v", err)
		res = &pb.HideQuestionRes{
			Code: http.StatusInternalServerError,
			Mag:  "internal err",
			Ok:   false,
		}
	}

	logger.Debugf("send message: %v", res.String())
	return res, nil
}
