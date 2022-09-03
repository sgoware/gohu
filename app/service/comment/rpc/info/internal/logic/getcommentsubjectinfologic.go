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

type GetCommentSubjectInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentSubjectInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentSubjectInfoLogic {
	return &GetCommentSubjectInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentSubjectInfoLogic) GetCommentSubjectInfo(in *pb.GetCommentSubjectInfoReq) (res *pb.GetCommentSubjectInfoRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	resData := &pb.GetCommentSubjectInfoRes_Data{}

	commentSubjectBytes, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("comment_subject_%d", in.SubjectId)).Bytes()
	if err != nil {
		logger.Errorf("get comment_subject cache failed, err: %v", err)

		commentSubjectModel := l.svcCtx.CommentModel.CommentSubject

		commentSubject, err := commentSubjectModel.WithContext(l.ctx).
			Where(commentSubjectModel.ID.Eq(in.SubjectId)).
			First()
		if err != nil {
			logger.Errorf("get commentSubject failed, err: mysql err, %v", err)
			res = &pb.GetCommentSubjectInfoRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}

		commentSubjectProto := &pb.CommentSubject{
			Id:         commentSubject.ID,
			ObjType:    commentSubject.ObjType,
			ObjId:      commentSubject.ObjID,
			Count:      commentSubject.Count,
			RootCount:  commentSubject.RootCount,
			State:      commentSubject.State,
			Attrs:      commentSubject.Attrs,
			CreateTime: commentSubject.CreateTime.String(),
			UpdateTime: commentSubject.UpdateTime.String(),
		}

		resData.CommentSubject = commentSubjectProto

		commentSubjectBytes, err = proto.Marshal(commentSubjectProto)
		if err != nil {
			logger.Errorf("marshal proto failed, err: %v", err)
		} else {
			l.svcCtx.Rdb.Set(l.ctx,
				fmt.Sprintf("comment_subject_%v", commentSubject.ID),
				commentSubjectBytes,
				time.Second*86400)
		}
	} else {
		commentSubject := &pb.CommentSubject{}
		err = proto.Unmarshal(commentSubjectBytes, commentSubject)
		if err != nil {
			logger.Errorf("unmarshal proto failed, err: %v", err)
			res = &pb.GetCommentSubjectInfoRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
		resData.CommentSubject = commentSubject
	}

	res = &pb.GetCommentSubjectInfoRes{
		Code: http.StatusOK,
		Msg:  "get comment subject successfully",
		Ok:   true,
		Data: resData,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
