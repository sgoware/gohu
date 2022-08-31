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

	commentIndex := &model.CommentIndex{}

	if in.RootId == 0 {
		// 是评论的情况
		count, err := commentIndexModel.WithContext(l.ctx).
			Where(commentIndexModel.SubjectID.Eq(in.SubjectId),
				commentIndexModel.RootID.Eq(0)).
			Count()
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

		commentIndex, err = commentIndexModel.WithContext(l.ctx).
			Where(commentIndexModel.SubjectID.Eq(in.SubjectId),
				commentIndexModel.UserID.Eq(j.Get("user_id").Int()),
				commentIndexModel.IPLoc.Eq(ip.GetIpLocFromApi(j.Get("last_ip").String())),
				commentIndexModel.CommentFloor.Eq(int32(count+1))).
			FirstOrCreate()
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
	} else {
		// 是回复评论的情况
		count, err := commentIndexModel.WithContext(l.ctx).
			Where(commentIndexModel.SubjectID.Eq(in.SubjectId),
				commentIndexModel.RootID.Eq(in.RootId)).
			Count()
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

		commentIndex, err = commentIndexModel.WithContext(l.ctx).
			Where(commentIndexModel.SubjectID.Eq(in.SubjectId),
				commentIndexModel.UserID.Eq(j.Get("user_id").Int()),
				commentIndexModel.IPLoc.Eq(ip.GetIpLocFromApi(j.Get("last_ip").String())),
				commentIndexModel.RootID.Eq(in.RootId),
				commentIndexModel.CommentID.Eq(in.CommentId),
				commentIndexModel.ReplyFloor.Eq(int32(count+1))).
			FirstOrCreate()
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
	}

	err = commentContentModel.WithContext(l.ctx).
		Create(&model.CommentContent{
			CommentID: commentIndex.ID,
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
