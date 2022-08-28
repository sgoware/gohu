package logic

import (
	"context"
	"gorm.io/gorm"
	"main/app/common/log"
	"main/app/service/question/rpc/crud/internal/svc"
	"main/app/service/question/rpc/crud/pb"
	"net/http"

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

	questionSubject, err := questionSubjectModel.WithContext(l.ctx).
		Select(questionSubjectModel.ID, questionSubjectModel.UserID, questionSubjectModel.State).
		Where(questionSubjectModel.ID.Eq(in.QuestionId)).First()
	switch err {
	case gorm.ErrRecordNotFound:
		res = &pb.HideQuestionRes{
			Code: http.StatusBadRequest,
			Msg:  "question not found",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil

	case nil:

	default:
		logger.Errorf("update question failed, err: %v", err)
		res = &pb.HideQuestionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
	}

	var state int32
	if questionSubject.State == 0 {
		state = 1
	} else {
		state = 0
	}

	_, err = questionSubjectModel.WithContext(l.ctx).Select(questionSubjectModel.State).
		Where(questionSubjectModel.ID.Eq(in.QuestionId)).
		Update(questionSubjectModel.State, state)
	if err != nil {
		logger.Errorf("update question failed, err: %v", err)
		res = &pb.HideQuestionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	res = &pb.HideQuestionRes{
		Code: http.StatusOK,
		Msg:  "hide question successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
