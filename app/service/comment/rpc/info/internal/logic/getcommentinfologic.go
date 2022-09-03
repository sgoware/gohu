package logic

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
	"main/app/common/log"
	"net/http"
	"time"

	"main/app/service/comment/rpc/info/internal/svc"
	"main/app/service/comment/rpc/info/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentInfoLogic {
	return &GetCommentInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentInfoLogic) GetCommentInfo(in *pb.GetCommentInfoReq) (res *pb.GetCommentInfoRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	resData := &pb.GetCommentInfoRes_Data{}

	commentIndexBytes, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("comment_index_%d", in.IndexId)).Bytes()
	if err != nil {
		logger.Errorf("get commentIndex cache failed, err: %v")

		commentIndexModel := l.svcCtx.CommentModel.CommentIndex

		commentIndex, err := commentIndexModel.WithContext(l.ctx).
			Where(commentIndexModel.ID.Eq(in.IndexId)).
			First()
		if err != nil {
			logger.Errorf("get commentIndex failed, err: mysql err, %v", err)
			res = &pb.GetCommentInfoRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		commentIndexProto := &pb.CommentIndex{
			Id:           commentIndex.ID,
			SubjectId:    commentIndex.SubjectID,
			UserId:       commentIndex.UserID,
			IpLoc:        commentIndex.IPLoc,
			RootId:       commentIndex.RootID,
			CommentFloor: commentIndex.CommentFloor,
			CommentId:    commentIndex.CommentID,
			ReplyFloor:   commentIndex.ReplyFloor,
			ApproveCount: commentIndex.ApproveCount,
			State:        commentIndex.State,
			Attrs:        commentIndex.Attrs,
			CreateTime:   commentIndex.CreateTime.String(),
			UpdateTime:   commentIndex.UpdateTime.String(),
		}

		resData.CommentIndex = commentIndexProto

		commentIndexBytes, err = proto.Marshal(commentIndexProto)
		if err != nil {
			logger.Errorf("marshal proto failed, err: %v", err)
		} else {
			l.svcCtx.Rdb.Set(l.ctx,
				fmt.Sprintf("comment_index_%v", commentIndex.ID),
				commentIndexBytes,
				time.Second*86400)
		}
	} else {
		commentIndex := &pb.CommentIndex{}
		err = proto.Unmarshal(commentIndexBytes, commentIndex)
		if err != nil {
			logger.Errorf("unmarshal proto failed, err: %v", err)
			res = &pb.GetCommentInfoRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		approveCnt, err := l.svcCtx.Rdb.Get(l.ctx,
			fmt.Sprintf("comment_index_approve_cnt_%d", in.IndexId)).Int()
		if err != nil {
			if err != redis.Nil {
				logger.Errorf("get [comment_subject_root_comment_cnt] failed, err: %v", err)
			}
		} else {
			resData.CommentIndex.ApproveCount += int32(approveCnt)
		}

		resData.CommentIndex = commentIndex
	}

	commentContentBytes, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("comment_content_%d", in.IndexId)).Bytes()
	if err != nil {
		logger.Errorf("get [comment_content] cache failed, err: %v", err)

		commentContentModel := l.svcCtx.CommentModel.CommentContent

		commentContent, err := commentContentModel.WithContext(l.ctx).
			Where(commentContentModel.CommentID.Eq(in.IndexId)).
			First()
		if err != nil {
			logger.Errorf("get [comment content] record failed, err: %v", err)
			res = &pb.GetCommentInfoRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		commentContentProto := &pb.CommentContent{
			CommentId:  commentContent.CommentID,
			Content:    commentContent.Content,
			Meta:       commentContent.Meta,
			CreateTime: commentContent.CreateTime.String(),
			UpdateTime: commentContent.UpdateTime.String(),
		}

		res.Data.CommentContent = commentContentProto

		commentContentBytes, err = proto.Marshal(commentContentProto)
		if err != nil {
			logger.Errorf("marshal proto failed, err: %v", err)
		} else {
			l.svcCtx.Rdb.Set(l.ctx,
				fmt.Sprintf("comment_content_%d", commentContent.CommentID),
				commentContentBytes,
				time.Second*86400)
		}
	} else {
		commentContent := &pb.CommentContent{}
		err = proto.Unmarshal(commentContentBytes, commentContent)
		if err != nil {
			logger.Errorf("unmarshal proto failed, err: %v", err)
			res = &pb.GetCommentInfoRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", err)
			return res, nil
		}
		resData.CommentContent = commentContent
	}

	res = &pb.GetCommentInfoRes{
		Code: http.StatusOK,
		Msg:  "get comment successfully",
		Ok:   true,
		Data: resData,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
