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
	res, _ := l.svcCtx.InfoRpcClient.GetQuestion(l.ctx, &info.GetQuestionReq{QuestionId: req.QuestionId})
	return &types.GetQuestionRes{
		Code: res.Code,
		Msg:  res.Msg,
		Ok:   res.Ok,
		Data: types.GetQuestionResData{
			QuestionSubject: types.QuestionSubject{
				Id:          res.Data.QuestionSubject.Id,
				UserId:      res.Data.QuestionSubject.UserId,
				IpLoc:       res.Data.QuestionSubject.IpLoc,
				Title:       res.Data.QuestionSubject.Title,
				Topic:       res.Data.QuestionSubject.Topic,
				Tag:         res.Data.QuestionSubject.Tag,
				SubCount:    res.Data.QuestionSubject.SubCount,
				AnswerCount: res.Data.QuestionSubject.AnswerCount,
				ViewCount:   res.Data.QuestionSubject.ViewCount,
				State:       res.Data.QuestionSubject.State,
				Attr:        res.Data.QuestionSubject.Attr,
				CreateTime:  res.Data.QuestionSubject.CreateTime,
				UpdateTime:  res.Data.QuestionSubject.UpdateTime,
			},
			QuestionContent: types.QuestionContent{
				QuestionId: res.Data.QuestionContent.QuestionId,
				Content:    res.Data.QuestionContent.Content,
				Meta:       res.Data.QuestionContent.Meta,
				CreateTime: res.Data.QuestionContent.CreateTime,
				UpdateTime: res.Data.QuestionContent.UpdateTime,
			},
		},
	}, nil
}
