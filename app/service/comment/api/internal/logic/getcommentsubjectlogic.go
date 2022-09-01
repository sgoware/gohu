package logic

import (
	"context"

	"main/app/service/comment/api/internal/svc"
	"main/app/service/comment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentSubjectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCommentSubjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentSubjectLogic {
	return &GetCommentSubjectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentSubjectLogic) GetCommentSubject(req *types.GetCommentSubjectReq) (resp *types.GetCommentSubjectRes, err error) {
	// todo: add your logic here and delete this line

	return
}
