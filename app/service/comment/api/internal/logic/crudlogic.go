package logic

import (
	"context"
	"encoding/json"
	"github.com/spf13/cast"
	"main/app/common/log"
	"main/app/service/comment/rpc/crud/crud"
	"net/http"

	"main/app/service/comment/api/internal/svc"
	"main/app/service/comment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CrudLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCrudLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CrudLogic {
	return &CrudLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CrudLogic) Crud(req *types.CrudReq) (resp *types.CrudRes, err error) {
	logger := log.GetSugaredLogger()
	res := &types.CrudRes{}

	if req.Action == "" || req.Data == "" {
		res = &types.CrudRes{
			Code: http.StatusBadRequest,
			Msg:  "param cannot be null",
			Ok:   false,
		}
		return res, nil
	}

	switch req.Action {
	case "publish":
		rpcReq := &crud.PublishCommentReq{
			UserDetails: "",
			SubjectId:   0,
			RootId:      0,
			CommentId:   0,
			Content:     "",
		}
		err = json.Unmarshal([]byte(req.Data), &rpcReq)
		logger.Debugf("rpcReq: %v", rpcReq)
		if err != nil {
			logger.Errorf("unmarshal data failed, err: %v", err)
			res = &types.CrudRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			return res, err
		}

		rpcReq.UserDetails = cast.ToString(l.ctx.Value("user_details"))
		rpcRes, _ := l.svcCtx.CrudRpcClient.PublishComment(l.ctx, rpcReq)
		res = &types.CrudRes{
			Code: rpcRes.Code,
			Msg:  rpcRes.Msg,
			Ok:   rpcRes.Ok,
		}
		return res, nil

	case "update":
		rpcReq := &crud.UpdateCommentReq{}
		err = json.Unmarshal([]byte(req.Data), &rpcReq)
		logger.Debugf("rpcReq: %v", rpcReq)
		if err != nil {
			logger.Errorf("unmarshal data failed, err: %v", err)
			res = &types.CrudRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			return res, err
		}

		rpcRes, _ := l.svcCtx.CrudRpcClient.UpdateComment(l.ctx, rpcReq)
		res = &types.CrudRes{
			Code: rpcRes.Code,
			Msg:  rpcRes.Msg,
			Ok:   rpcRes.Ok,
		}
		return res, nil

	case "delete":
		rpcReq := &crud.DeleteCommentReq{}
		err = json.Unmarshal([]byte(req.Data), &rpcReq)
		logger.Debugf("rpcReq: %v", rpcReq)
		if err != nil {
			logger.Errorf("unmarshal data failed, err: %v", err)
			res = &types.CrudRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			return res, err
		}

		rpcRes, _ := l.svcCtx.CrudRpcClient.DeleteComment(l.ctx, rpcReq)
		res = &types.CrudRes{
			Code: rpcRes.Code,
			Msg:  rpcRes.Msg,
			Ok:   rpcRes.Ok,
		}
		return res, nil

	default:
		res = &types.CrudRes{
			Code: http.StatusBadRequest,
			Msg:  `param "object" err`,
			Ok:   false,
		}
		return res, nil
	}
}
