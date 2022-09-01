package logic

import (
	"context"
	"main/app/common/log"
	"main/app/service/comment/rpc/info/internal/svc"
	"main/app/service/comment/rpc/info/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentSubjectIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentSubjectIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentSubjectIndexLogic {
	return &GetCommentSubjectIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentSubjectIndexLogic) GetCommentSubjectIndex(in *pb.GetCommentSubjectIndexReq) (res *pb.GetCommentSubjectIndexRes, err error) {
	// todo: add your logic here and delete this line
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	logger.Debugf("send message: %v", res.String())
	return res, nil
}
