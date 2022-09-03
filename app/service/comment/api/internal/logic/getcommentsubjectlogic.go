package logic

import (
	"context"
	"main/app/service/comment/rpc/info/info"

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
	rpcRes, _ := l.svcCtx.InfoRpcClient.GetCommentSubject(l.ctx, &info.GetCommentSubjectReq{SubjectId: req.SubjectId})

	return &types.GetCommentSubjectRes{
		Code: rpcRes.Code,
		Msg:  rpcRes.Msg,
		Ok:   rpcRes.Ok,
		Data: types.GetCommentSubjectResData{CommentSubject: types.CommentSubject{
			Id:         rpcRes.Data.CommentSubject.Id,
			ObjType:    rpcRes.Data.CommentSubject.ObjType,
			ObjId:      rpcRes.Data.CommentSubject.ObjId,
			Count:      rpcRes.Data.CommentSubject.Count,
			RootCount:  rpcRes.Data.CommentSubject.Count,
			State:      rpcRes.Data.CommentSubject.State,
			Attrs:      rpcRes.Data.CommentSubject.Attrs,
			CreateTime: rpcRes.Data.CommentSubject.CreateTime,
			UpdateTime: rpcRes.Data.CommentSubject.UpdateTime,
		}},
	}, nil
}
