package logic

import (
	"context"

	"main/app/service/comment/api/internal/svc"
	"main/app/service/comment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentSubjectIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCommentSubjectIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentSubjectIndexLogic {
	return &GetCommentSubjectIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentSubjectIndexLogic) GetCommentSubjectIndex(req *types.GetCommenSubjectIndexReq) (resp *types.GetCommenSubjectIndexRes, err error) {
	// todo: add your logic here and delete this line

	return
}
