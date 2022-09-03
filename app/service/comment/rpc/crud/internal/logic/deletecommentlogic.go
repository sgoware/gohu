package logic

import (
	"context"
	"fmt"
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

	// 级联删除commentContent
	commentIndexModel := l.svcCtx.CommentModel.CommentIndex

	commentIndex, err := commentIndexModel.WithContext(l.ctx).
		Select(commentIndexModel.CommentID, commentIndexModel.SubjectID, commentIndexModel.RootID).
		Where(commentIndexModel.CommentID.Eq(in.CommentId)).
		First()
	if err != nil {
		logger.Errorf("query [comment_index] record failed, err: %v", err)
		res = &pb.DeleteCommentRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	if commentIndex.RootID == 0 {
		err = l.svcCtx.Rdb.Decr(l.ctx,
			fmt.Sprintf("comment_subject_root_comment_cnt_%d", commentIndex.SubjectID)).Err()
		if err != nil {
			logger.Errorf("incr [comment_subject_root_comment_cnt] failed, err: %v", err)
			res = &pb.DeleteCommentRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		err = l.svcCtx.Rdb.SAdd(l.ctx,
			"comment_subject_root_comment_cnt_set",
			commentIndex.SubjectID).Err()
		if err != nil {
			logger.Errorf("update [comment_subject_root_comment_cnt_set] failed, err: %v", err)
			res = &pb.DeleteCommentRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
	}

	err = l.svcCtx.Rdb.Decr(l.ctx,
		fmt.Sprintf("comment_subject_comment_cnt_%d", commentIndex.SubjectID)).Err()
	if err != nil {
		logger.Errorf("incr [comment_subject_comment_cnt] failed, err: %v", err)
		res = &pb.DeleteCommentRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	err = l.svcCtx.Rdb.SAdd(l.ctx,
		"comment_subject_comment_cnt_set",
		in.CommentId).Err()
	if err != nil {
		logger.Errorf("update [comment_subject_comment_cnt_set] failed, err: %v", err)
		res = &pb.DeleteCommentRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	_, err = commentIndexModel.WithContext(l.ctx).
		Where(commentIndexModel.ID.Eq(in.CommentId)).
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

	err = l.svcCtx.Rdb.SRem(l.ctx,
		fmt.Sprintf("comment_id_user_set_%v", commentIndex.UserID),
		commentIndex.ID).Err()
	if err != nil {
		logger.Errorf("update [comment_id_user_set] failed, err: %v", err)
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
