package logic

import (
	"context"
	"main/app/common/log"
	"net/http"

	"main/app/service/comment/rpc/crud/internal/svc"
	"main/app/service/comment/rpc/crud/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCommentLogic) DeleteComment(in *pb.DeleteCommentReq) (res *pb.DeleteCommentRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	commentIndexModel := l.svcCtx.CommentModel.CommentIndex
	commentContentModel := l.svcCtx.CommentModel.CommentContent

	_, err = commentIndexModel.WithContext(l.ctx).
		Where(commentIndexModel.CommentID.Eq(in.CommentId)).
		Delete()
	if err != nil {
		logger.Errorf("delete comment failed, err: mysql err, %v", err)
		res = &pb.DeleteCommentRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	_, err = commentContentModel.WithContext(l.ctx).
		Where(commentContentModel.CommentID.Eq(in.CommentId)).
		Delete()
	if err != nil {
		logger.Errorf("delete comment failed, err: mysql err, %v", err)
		res = &pb.DeleteCommentRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	res = &pb.DeleteCommentRes{
		Code: http.StatusOK,
		Msg:  "delete comment successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
