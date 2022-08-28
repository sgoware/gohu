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

	questionContentModel := l.svcCtx.QuestionModel.QuestionContent

	_, err = questionContentModel.WithContext(l.ctx).Select(
		questionContentModel.Title,
		questionContentModel.Content,
		questionContentModel.Topic,
		questionContentModel.Tag).
		Where(questionContentModel.QuestionID.Eq(in.QuestionId)).
		Updates(model.QuestionContent{
			Title:   in.Title,
			Topic:   in.Topic,
			Tag:     in.Tag,
			Content: in.Content,
		})
	switch err {
	case gorm.ErrRecordNotFound:
		res = &pb.UpdateQuestionRes{
			Code: http.StatusBadRequest,
			Mag:  "question not found",
			Ok:   false,
		}
	case nil:
		res = &pb.UpdateQuestionRes{
			Code: http.StatusOK,
			Mag:  "update question successfully",
			Ok:   true,
		}
	default:
		logger.Errorf("update question failed, err: %v", err)
		res = &pb.UpdateQuestionRes{
			Code: http.StatusInternalServerError,
			Mag:  "internal err",
			Ok:   false,
		}
	}

	logger.Debugf("send message: %v", res.String())
	return res, nil
}
