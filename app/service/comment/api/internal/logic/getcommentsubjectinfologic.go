package logic

import (
	"context"
	"main/app/service/comment/rpc/info/info"

	"main/app/service/comment/api/internal/svc"
	"main/app/service/comment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentSubjectInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCommentSubjectInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentSubjectInfoLogic {
	return &GetCommentSubjectInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentSubjectInfoLogic) GetCommentSubjectInfo(req *types.GetCommentSubjectInfoReq) (resp *types.GetCommentSubjectInfoRes, err error) {
	rpcRes, _ := l.svcCtx.InfoRpcClient.GetCommentSubjectInfo(l.ctx, &info.GetCommentSubjectInfoReq{SubjectId: req.SubjectId})

	return &types.GetCommentSubjectInfoRes{
		Code: rpcRes.Code,
		Msg:  rpcRes.Msg,
		Ok:   rpcRes.Ok,
		Data: types.GetCommentSubjectInfoResData{CommentSubject: types.CommentSubject{
			Id:         rpcRes.Data.CommentSubject.Id,
			ObjType:    rpcRes.Data.CommentSubject.ObjType,
			ObjId:      rpcRes.Data.CommentSubject.ObjId,
			Count:      rpcRes.Data.CommentSubject.Count,
			RootCount:  rpcRes.Data.CommentSubject.RootCount,
			State:      rpcRes.Data.CommentSubject.State,
			Attrs:      rpcRes.Data.CommentSubject.Attrs,
			CreateTime: rpcRes.Data.CommentSubject.CreateTime,
			UpdateTime: rpcRes.Data.CommentSubject.UpdateTime,
		}},
	}, nil
}
