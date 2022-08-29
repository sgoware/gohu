package logic

import (
	"context"
	"main/app/common/log"
	"net/http"

	"main/app/service/comment/rpc/crud/internal/svc"
	"main/app/service/comment/rpc/crud/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitSubjectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitSubjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitSubjectLogic {
	return &InitSubjectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InitSubjectLogic) InitSubject(in *pb.InitSubjectReq) (res *pb.InitSubjectRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	commentSubjectModel := l.svcCtx.CommentModel.CommentSubject

	_, err = commentSubjectModel.WithContext(l.ctx).
		Where(commentSubjectModel.ObjType.Eq(in.ObjType), commentSubjectModel.ID.Eq(in.ObjId)).
		FirstOrCreate()
	if err != nil {
		logger.Errorf("init comment subject failed, err: mysql err, %v", err)
		res = &pb.InitSubjectRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	res = &pb.InitSubjectRes{
		Code: http.StatusOK,
		Msg:  "init comment subject successfully",
		Ok:   true,
	}
	logger.Debugf("send message: %v", res.String())
	return res, nil
}
