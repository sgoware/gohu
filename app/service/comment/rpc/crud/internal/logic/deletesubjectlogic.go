package logic

import (
	"context"
	"gorm.io/gorm"
	"main/app/common/log"
	"net/http"

	"main/app/service/comment/rpc/crud/internal/svc"
	"main/app/service/comment/rpc/crud/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSubjectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSubjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSubjectLogic {
	return &DeleteSubjectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteSubjectLogic) DeleteSubject(in *pb.DeleteSubjectReq) (res *pb.DeleteSubjectRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	commentSubjectModel := l.svcCtx.CommentModel.CommentSubject

	_, err = commentSubjectModel.WithContext(l.ctx).
		Select(commentSubjectModel.ObjType, commentSubjectModel.ObjID).
		Where(commentSubjectModel.ObjType.Eq(in.ObjType), commentSubjectModel.ObjID.Eq(in.ObjId)).
		Delete()
	switch err {
	case gorm.ErrRecordNotFound:
		res = &pb.DeleteSubjectRes{
			Code: http.StatusForbidden,
			Msg:  "comment subject not found",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil

	case nil:
		res = &pb.DeleteSubjectRes{
			Code: http.StatusOK,
			Msg:  "delete comment subject successfully",
			Ok:   true,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil

	default:
		logger.Errorf("delete comment subject failed, err: mysql err, %v", err)
		res = &pb.DeleteSubjectRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		return res, nil
	}
}
