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

type DeleteQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteQuestionLogic {
	return &DeleteQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteQuestionLogic) DeleteQuestion(in *pb.DeleteQuestionReq) (res *pb.DeleteQuestionRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	// 删除question_subject后, 级联删除关联的回答和回答下的评论
	questionSubjectModel := l.svcCtx.QuestionModel.QuestionSubject

	_, err = questionSubjectModel.WithContext(l.ctx).Where(questionSubjectModel.ID.Eq(in.QuestionId)).Delete()
	switch err {
	case gorm.ErrRecordNotFound:
		res = &pb.DeleteQuestionRes{
			Code: http.StatusBadRequest,
			Msg:  "question not found",
			Ok:   false,
		}
	case nil:

	default:
		logger.Errorf("update question failed, err: %v", err)
		res = &pb.DeleteQuestionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
	}

	res = &pb.DeleteQuestionRes{
		Code: http.StatusOK,
		Msg:  "delete question successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
