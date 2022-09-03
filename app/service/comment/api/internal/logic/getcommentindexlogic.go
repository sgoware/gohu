package logic

import (
	"context"
	"main/app/service/comment/rpc/info/info"

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
	rpcRes, _ := l.svcCtx.InfoRpcClient.GetCommentInfo(l.ctx, &info.GetCommentInfoReq{IndexId: req.IndexId})

	return &types.GetCommentIndexRes{
		Code: rpcRes.Code,
		Msg:  rpcRes.Msg,
		Ok:   rpcRes.Ok,
		Data: types.GetCommentIndexResData{
			CommentIndex: types.CommentIndex{
				Id:           rpcRes.Data.CommentIndex.Id,
				SubjectId:    rpcRes.Data.CommentIndex.SubjectId,
				UserId:       rpcRes.Data.CommentIndex.UserId,
				IpLoc:        rpcRes.Data.CommentIndex.IpLoc,
				RootId:       rpcRes.Data.CommentIndex.RootId,
				CommentFloor: rpcRes.Data.CommentIndex.CommentFloor,
				CommentId:    rpcRes.Data.CommentIndex.CommentId,
				ReplyFloor:   rpcRes.Data.CommentIndex.ReplyFloor,
				ApproveCount: rpcRes.Data.CommentIndex.ApproveCount,
				State:        rpcRes.Data.CommentIndex.State,
				Attrs:        rpcRes.Data.CommentIndex.Attrs,
				CreateTime:   rpcRes.Data.CommentIndex.CreateTime,
				UpdateTime:   rpcRes.Data.CommentIndex.UpdateTime,
			},
			CommentContent: types.CommentContent{
				CommentId:  rpcRes.Data.CommentContent.CommentId,
				Content:    rpcRes.Data.CommentContent.Content,
				Meta:       rpcRes.Data.CommentContent.Meta,
				CreateTime: rpcRes.Data.CommentContent.CreateTime,
				UpdateTime: rpcRes.Data.CommentContent.UpdateTime,
			},
		},
	}, nil
}
