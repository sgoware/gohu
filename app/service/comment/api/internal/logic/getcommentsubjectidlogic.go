package logic

import (
	"context"
	"main/app/service/comment/rpc/info/info"

	"main/app/service/comment/api/internal/svc"
	"main/app/service/comment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentSubjectIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCommentSubjectIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentSubjectIdLogic {
	return &GetCommentSubjectIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentSubjectIdLogic) GetCommentSubjectId(req *types.GetCommentSubjectIdReq) (resp *types.GetCommentSubjectIdRes, err error) {
	rpcRes, _ := l.svcCtx.InfoRpcClient.GetCommentSubjectId(l.ctx, &info.GetCommentSubjectIdReq{
		ObjType: req.ObjType,
		ObjId:   req.ObjId,
	})

	return &types.GetCommentSubjectIdRes{
		Code: rpcRes.Code,
		Msg:  rpcRes.Msg,
		Ok:   rpcRes.Ok,
		Data: types.GetCommentSubjectIdResData{SubjectId: rpcRes.Data.SubjectId},
	}, nil
}
