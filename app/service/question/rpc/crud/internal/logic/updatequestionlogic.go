package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"main/app/common/log"
	"main/app/service/question/dao/model"
	"main/app/service/question/rpc/crud/internal/svc"
	"main/app/service/question/rpc/crud/pb"
	"net/http"
)

type UpdateQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateQuestionLogic {
	return &UpdateQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateQuestionLogic) UpdateQuestion(in *pb.UpdateQuestionReq) (res *pb.UpdateQuestionRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	questionSubjectModel := l.svcCtx.QuestionModel.QuestionSubject
	questionContentModel := l.svcCtx.QuestionModel.QuestionContent

	_, err = questionSubjectModel.WithContext(l.ctx).Select(
		questionSubjectModel.ID,
		questionSubjectModel.Title,
		questionSubjectModel.Topic,
		questionSubjectModel.Tag,
	).
		Where(questionSubjectModel.ID.Eq(in.QuestionId)).
		Updates(model.QuestionSubject{
			Title: in.Title,
			Topic: in.Topic,
			Tag:   in.Tag,
		})

	_, err = questionContentModel.WithContext(l.ctx).Select(
		questionContentModel.QuestionID,
		questionContentModel.Content,
	).
		Where(questionContentModel.QuestionID.Eq(in.QuestionId)).
		Update(questionContentModel.Content, in.Content)
	switch err {
	case gorm.ErrRecordNotFound:
		res = &pb.UpdateQuestionRes{
			Code: http.StatusBadRequest,
			Msg:  "question not found",
			Ok:   false,
		}
	case nil:
		res = &pb.UpdateQuestionRes{
			Code: http.StatusOK,
			Msg:  "update question successfully",
			Ok:   true,
		}
	default:
		logger.Errorf("update question failed, err: %v", err)
		res = &pb.UpdateQuestionRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
	}

	logger.Debugf("send message: %v", res.String())
	return res, nil
}
