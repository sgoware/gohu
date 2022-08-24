package logic

import (
	"context"
	"main/app/service/user/api/internal/svc"
	"main/app/service/user/api/internal/types"
	"main/app/service/user/rpc/crud/crud"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterRes, err error) {
	res, err := l.svcCtx.CrudRpcClient.Register(l.ctx, &crud.RegisterReq{
		Uid:      req.Uid,
		Nickname: req.Nickname,
		Password: req.Password,
	})
	if err != nil {
		logx.Errorf("create user failed, err: %v", err)
	}
	return &types.RegisterRes{
		Code: int(res.Code),
		Msg:  res.Msg,
	}, nil
}
