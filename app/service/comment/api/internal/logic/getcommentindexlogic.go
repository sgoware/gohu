package logic

import (
	"context"

	"main/app/service/comment/api/internal/svc"
	"main/app/service/comment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCommentIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentIndexLogic {
	return &GetCommentIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentIndexLogic) GetCommentIndex(req *types.GetCommentIndexReq) (resp *types.GetCommentIndexRes, err error) {
	// todo: add your logic here and delete this line

	return
}
