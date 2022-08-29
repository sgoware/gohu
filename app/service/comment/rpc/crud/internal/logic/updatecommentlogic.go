package logic

import (
	"context"
	"main/app/common/log"
	"net/http"

	"main/app/service/comment/rpc/crud/internal/svc"
	"main/app/service/comment/rpc/crud/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCommentLogic {
	return &UpdateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCommentLogic) UpdateComment(in *pb.UpdateCommentReq) (res *pb.UpdateCommentRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	answerContentModel := l.svcCtx.CommentModel.CommentContent

	_, err = answerContentModel.WithContext(l.ctx).
		Select(answerContentModel.CommentID, answerContentModel.Content).
		Where(answerContentModel.CommentID.Eq(in.CommentId)).
		Update(answerContentModel.Content, in.Content)
	if err != nil {
		logger.Errorf("update comment failed, err: mysql err, %v", err)
		res = &pb.UpdateCommentRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	res = &pb.UpdateCommentRes{
		Code: http.StatusOK,
		Msg:  "update comment successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
