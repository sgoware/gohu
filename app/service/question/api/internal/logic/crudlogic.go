package logic

import (
	"context"
	"encoding/json"
	"github.com/spf13/cast"
	"main/app/common/log"
	"main/app/service/question/rpc/crud/crud"
	"net/http"

	"main/app/service/question/api/internal/svc"
	"main/app/service/question/api/internal/types"

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

func (l *CrudLogic) Crud(req *types.CrudReq) (res *types.CrudRes, err error) {
	logger := log.GetSugaredLogger()
	res = &types.CrudRes{}

	if req.Object == "" || req.Action == "" {
		res = &types.CrudRes{
			Code: http.StatusBadRequest,
			Msg:  "param cannot be null",
			Ok:   false,
		}
		return res, nil
	}
	switch req.Object {
	case "question":
		switch req.Action {
		case "publish":
			rpcReq := &crud.PublishQuestionReq{}
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
			rpcRes, err := l.svcCtx.CrudRpcClient.PublishQuestion(l.ctx, rpcReq)
			if err != nil {
				res = &types.CrudRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				return res, err
			}

			res = &types.CrudRes{
				Code: int(rpcRes.Code),
				Msg:  rpcRes.Msg,
				Ok:   rpcRes.Ok,
			}
			return res, err
			// TODO: 使用队列发布消息
			// l.svcCtx.

			return res, nil

		case "update":
			rpcReq := &crud.UpdateQuestionReq{}
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
			rpcRes, err := l.svcCtx.CrudRpcClient.UpdateQuestion(l.ctx, rpcReq)
			if err != nil {
				res = &types.CrudRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				return res, err
			}

			res = &types.CrudRes{
				Code: int(rpcRes.Code),
				Msg:  rpcRes.Msg,
				Ok:   rpcRes.Ok,
			}
			return res, err

		case "hide":
			rpcReq := &crud.HideQuestionReq{}
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
			rpcRes, err := l.svcCtx.CrudRpcClient.HideQuestion(l.ctx, rpcReq)
			if err != nil {
				res = &types.CrudRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				return res, err
			}

			res = &types.CrudRes{
				Code: int(rpcRes.Code),
				Msg:  rpcRes.Msg,
				Ok:   rpcRes.Ok,
			}
			return res, err

		case "delete":
			rpcReq := &crud.DeleteQuestionReq{}
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
			rpcRes, err := l.svcCtx.CrudRpcClient.DeleteQuestion(l.ctx, rpcReq)
			if err != nil {
				logger.Errorf("mapping struct failed, err: %v", err)
				res = &types.CrudRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				return res, err
			}

			res = &types.CrudRes{
				Code: int(rpcRes.Code),
				Msg:  rpcRes.Msg,
				Ok:   rpcRes.Ok,
			}
			return res, err

		default:
			res = &types.CrudRes{
				Code: http.StatusBadRequest,
				Msg:  `param "action" err`,
				Ok:   false,
			}
			return res, nil
		}

	case "answer":
		switch req.Action {
		case "publish":
			rpcReq := &crud.PublishAnswerReq{}
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
			rpcRes, err := l.svcCtx.CrudRpcClient.PublishAnswer(l.ctx, rpcReq)
			if err != nil {
				res = &types.CrudRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				return res, err
			}

			res = &types.CrudRes{
				Code: int(rpcRes.Code),
				Msg:  rpcRes.Msg,
				Ok:   rpcRes.Ok,
			}
			return res, err

		case "update":
			rpcReq := &crud.UpdateAnswerReq{}
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
			rpcRes, err := l.svcCtx.CrudRpcClient.UpdateAnswer(l.ctx, rpcReq)
			if err != nil {
				res = &types.CrudRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				return res, err
			}

			res = &types.CrudRes{
				Code: int(rpcRes.Code),
				Msg:  rpcRes.Msg,
				Ok:   rpcRes.Ok,
			}
			return res, err

		case "hide":
			rpcReq := &crud.HideAnswerReq{}
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
			rpcRes, err := l.svcCtx.CrudRpcClient.HideAnswer(l.ctx, rpcReq)
			if err != nil {
				res = &types.CrudRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				return res, err
			}

			res = &types.CrudRes{
				Code: int(rpcRes.Code),
				Msg:  rpcRes.Msg,
				Ok:   rpcRes.Ok,
			}
			return res, err

		case "delete":
			rpcReq := &crud.DeleteAnswerReq{}
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
			rpcRes, err := l.svcCtx.CrudRpcClient.DeleteAnswer(l.ctx, rpcReq)
			if err != nil {
				res = &types.CrudRes{
					Code: http.StatusInternalServerError,
					Msg:  "internal err",
					Ok:   false,
				}
				return res, err
			}

			res = &types.CrudRes{
				Code: int(rpcRes.Code),
				Msg:  rpcRes.Msg,
				Ok:   rpcRes.Ok,
			}
			return res, err

		default:
			res = &types.CrudRes{
				Code: http.StatusBadRequest,
				Msg:  `param "action" err`,
				Ok:   false,
			}
			return res, nil
		}

	default:
		res = &types.CrudRes{
			Code: http.StatusBadRequest,
			Msg:  `param "object" err`,
			Ok:   false,
		}
		return res, nil
	}

	return
}
