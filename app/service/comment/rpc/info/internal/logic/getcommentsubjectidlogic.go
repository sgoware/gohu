package logic

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"main/app/common/log"
	"net/http"
	"time"

	"main/app/service/comment/rpc/info/internal/svc"
	"main/app/service/comment/rpc/info/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentSubjectIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentSubjectIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentSubjectIdLogic {
	return &GetCommentSubjectIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentSubjectIdLogic) GetCommentSubjectId(in *pb.GetCommentSubjectIdReq) (res *pb.GetCommentSubjectIdRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	id, err := l.svcCtx.Rdb.Get(l.ctx,
		fmt.Sprintf("comment_subject_id_%d_%d", in.ObjType, in.ObjId)).Result()
	if err == nil {
		res = &pb.GetCommentSubjectIdRes{
			Code: http.StatusOK,
			Msg:  "get comment subject id successfully",
			Ok:   true,
			Data: &pb.GetCommentSubjectIdRes_Data{SubjectId: cast.ToInt64(id)},
		}
	}
	commentSubjectModel := l.svcCtx.CommentModel.CommentSubject
	commentSubject, err := commentSubjectModel.WithContext(l.ctx).
		Select(commentSubjectModel.ID, commentSubjectModel.ObjType, commentSubjectModel.ObjID).
		Where(commentSubjectModel.ObjType.Eq(in.ObjType), commentSubjectModel.ObjID.Eq(in.ObjId)).
		First()
	if err != nil {
		logger.Errorf("query [comment_subject] record failed, err: %v", err)
		res = &pb.GetCommentSubjectIdRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	err = l.svcCtx.Rdb.Set(l.ctx,
		fmt.Sprintf("comment_subject_id_%d_%d", in.ObjType, in.ObjId),
		commentSubject.ID,
		time.Second*86400).Err()
	if err != nil {
		logger.Errorf("set [comment_subject_id] cache failed, err: %v", err)
	}

	res = &pb.GetCommentSubjectIdRes{
		Code: http.StatusOK,
		Msg:  "get comment subject id successfully",
		Ok:   true,
		Data: &pb.GetCommentSubjectIdRes_Data{SubjectId: commentSubject.ID},
	}

	logger.Debugf("send message: %v", res.String())
	return res, nil
}
