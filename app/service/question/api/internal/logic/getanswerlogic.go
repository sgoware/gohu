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
		Data: res.Data.String(), // TODO: 待测试
	}, nil
}
