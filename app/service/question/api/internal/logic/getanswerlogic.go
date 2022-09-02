package logic

import (
	"context"
	"main/app/service/question/rpc/info/info"

	"main/app/service/question/api/internal/svc"
	"main/app/service/question/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAnswerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAnswerLogic {
	return &GetAnswerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAnswerLogic) GetAnswer(req *types.GetAnswerReq) (resp *types.GetAnswerRes, err error) {
	res, _ := l.svcCtx.InfoRpcClient.GetAnswer(l.ctx, &info.GetAnswerReq{AnswerId: req.Id})
	return &types.GetAnswerRes{
		Code: res.Code,
		Msg:  res.Msg,
		Ok:   res.Ok,
		Data: types.GetAnswerResData{
			AnswerIndex: types.AnswerIndex{
				Id:           res.Data.AnswerIndex.Id,
				QuestionId:   res.Data.AnswerIndex.QuestionId,
				UserId:       res.Data.AnswerIndex.UserId,
				IpLoc:        res.Data.AnswerIndex.IpLoc,
				ApproveCount: res.Data.AnswerIndex.ApproveCount,
				LikeCount:    res.Data.AnswerIndex.LikeCount,
				CollectCount: res.Data.AnswerIndex.CollectCount,
				State:        res.Data.AnswerIndex.State,
				Attrs:        res.Data.AnswerIndex.Attrs,
				CreateTime:   res.Data.AnswerIndex.CreateTime,
				UpdateTime:   res.Data.AnswerIndex.UpdateTime,
			},
			AnswerContent: types.AnswerContent{
				AnswerId:   res.Data.AnswerContent.AnswerId,
				Content:    res.Data.AnswerContent.Content,
				Meta:       res.Data.AnswerContent.Meta,
				CreateTime: res.Data.AnswerContent.CreateTime,
				UpdateTime: res.Data.AnswerContent.UpdateTime,
			},
		},
	}, nil
}
