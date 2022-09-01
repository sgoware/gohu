package logic

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"main/app/common/log"
	"net/http"
	"time"

	"main/app/service/comment/rpc/info/internal/svc"
	"main/app/service/comment/rpc/info/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentIndexLogic {
	return &GetCommentIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentIndexLogic) GetCommentIndex(in *pb.GetCommentIndexReq) (res *pb.GetCommentIndexRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	resData := &pb.GetCommentIndexRes_Data{}

	commentIndexBytes, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("commentIndex_%d", in.IndexId)).Bytes()
	if err != nil {
		logger.Errorf("get commentIndex cache failed, err: %v")

		commentIndexModel := l.svcCtx.CommentModel.CommentIndex

		commentIndex, err := commentIndexModel.WithContext(l.ctx).
			Where(commentIndexModel.ID.Eq(in.IndexId)).
			First()
		if err != nil {
			logger.Errorf("get commentIndex failed, err: mysql err, %v", err)
			res = &pb.GetCommentIndexRes{
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
				fmt.Sprintf("commentIndex_%v", commentIndex.ID),
				commentIndexBytes,
				time.Second*86400)
		}
	} else {
		commentIndex := &pb.CommentIndex{}
		err = proto.Unmarshal(commentIndexBytes, commentIndex)
		if err != nil {
			logger.Errorf("unmarshal proto failed, err: %v", err)
			res = &pb.GetCommentIndexRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		resData.CommentIndex = commentIndex
	}

	res = &pb.GetCommentIndexRes{
		Code: http.StatusOK,
		Msg:  "get comment index successfully",
		Ok:   true,
		Data: resData,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
