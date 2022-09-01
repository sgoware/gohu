package logic

import (
	"context"
	"main/app/service/question/rpc/info/info"

	"main/app/service/question/api/internal/svc"
	"main/app/service/question/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionLogic {
	return &GetQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQuestionLogic) GetQuestion(req *types.GetQuestionReq) (resp *types.GetQuestionRes, err error) {
	res, _ := l.svcCtx.InfoRpcClient.GetQuestion(l.ctx, &info.GetQuestionReq{QuestionId: req.Id})
	return &types.GetQuestionRes{
		Code: res.Code,
		Msg:  res.Msg,
		Ok:   res.Ok,
		Data: res.Data.String(), // TODO: 待测试
	}, nil
}
