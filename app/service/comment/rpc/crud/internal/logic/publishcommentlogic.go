package logic

import (
	"context"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"main/app/common/log"
	"main/app/service/comment/dao/model"
	"main/app/service/comment/rpc/crud/internal/svc"
	"main/app/service/comment/rpc/crud/pb"
	"main/app/utils/net/ip"
	"net/http"
)

type PublishCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishCommentLogic {
	return &PublishCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishCommentLogic) PublishComment(in *pb.PublishCommentReq) (res *pb.PublishCommentRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	j := gjson.Parse(in.UserDetails)

	commentIndexModel := l.svcCtx.CommentModel.CommentIndex
	commentContentModel := l.svcCtx.CommentModel.CommentContent

	err = commentIndexModel.WithContext(l.ctx).
		Create(&model.CommentIndex{
			SubjectID:    in.SubjectId,
			UserID:       j.Get("user_id").Int(),
			IPLoc:        ip.GetIpLocFromApi(j.Get("last_ip").String()),
			RootID:       in.RootId,
			CommentFloor: in.CommentFloor,
			CommentID:    in.CommentId,
			ReplyFloor:   in.ReplyFloor,
		})
	if err != nil {
		logger.Errorf("publish comment failed, err: mysql err, %v", err)
		res = &pb.PublishCommentRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	commentIndex, err := commentIndexModel.WithContext(l.ctx).
		Select(commentIndexModel.ID, commentIndexModel.UserID).
		Where(commentIndexModel.UserID.Eq(j.Get("user_id").Int())).
		Order(commentIndexModel.ID.Desc()).Last()
	if err != nil {
		logger.Errorf("publish comment failed, err: mysql err, %v", err)
		res = &pb.PublishCommentRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	err = commentContentModel.WithContext(l.ctx).
		Create(&model.CommentContent{
			CommentID: commentIndex.CommentID,
			Content:   in.Content,
		})
	if err != nil {
		logger.Errorf("publish comment failed, err: mysql err, %v", err)
		res = &pb.PublishCommentRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	res = &pb.PublishCommentRes{
		Code: http.StatusOK,
		Msg:  "publish comment successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
