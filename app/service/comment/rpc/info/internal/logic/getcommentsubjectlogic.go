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

type GetCommentSubjectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentSubjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentSubjectLogic {
	return &GetCommentSubjectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentSubjectLogic) GetCommentSubject(in *pb.GetCommentSubjectReq) (res *pb.GetCommentSubjectRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	resData := &pb.GetCommentSubjectRes_Data{}

	commentSubjectBytes, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("commentSubject_%d", in.SubjectId)).Bytes()
	if err != nil {
		logger.Errorf("get commentSubject cache failed, err: %v", err)

		commentSubjectModel := l.svcCtx.CommentModel.CommentSubject

		commentSubject, err := commentSubjectModel.WithContext(l.ctx).
			Where(commentSubjectModel.ID.Eq(in.SubjectId)).
			First()
		if err != nil {
			logger.Errorf("get commentSubject failed, err: mysql err, %v", err)
			res = &pb.GetCommentSubjectRes{
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
				fmt.Sprintf("commentSubject_%v", commentSubject.ID),
				commentSubjectBytes,
				time.Second*86400)
		}
	} else {
		commentSubject := &pb.CommentSubject{}
		err = proto.Unmarshal(commentSubjectBytes, commentSubject)
		if err != nil {
			logger.Errorf("unmarshal proto failed, err: %v", err)
			res = &pb.GetCommentSubjectRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, nil
		}
		resData.CommentSubject = commentSubject
	}

	res = &pb.GetCommentSubjectRes{
		Code: http.StatusOK,
		Msg:  "get comment subject successfully",
		Ok:   true,
		Data: resData,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
