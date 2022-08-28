package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"main/app/common/log"
	"main/app/service/oauth/model"
	"main/app/service/oauth/rpc/token/store/internal/svc"
	"main/app/service/oauth/rpc/token/store/pb"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenLogic {
	return &GetTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTokenLogic) GetToken(in *pb.GetTokenReq) (res *pb.GetTokenRes, err error) {
	logger := log.GetSugaredLogger()
	logger.Debugf("recv message: %v", in.String())

	if in.UserId == 0 {
		res = &pb.GetTokenRes{
			Code: http.StatusBadRequest,
			Msg:  fmt.Sprintf("get token failed, %v", model.ErrInvalidUserId),
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	}

	val, err := l.svcCtx.Rdb.Get(l.ctx, model.JwtToken+"_"+strconv.FormatInt(in.UserId, 10)).Result()
	if err != nil {
		res = &pb.GetTokenRes{
			Code: http.StatusOK,
			Msg:  "token not found",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	} else {
		res = &pb.GetTokenRes{
			Code: http.StatusOK,
			Msg:  "get token successfully",
			Ok:   true,
			Data: &pb.GetTokenRes_Data{OauthToken: &pb.OAuth2Token{}},
		}
		err = json.Unmarshal([]byte(val), &res.Data.OauthToken)
		if err != nil {
			logger.Errorf("unmarshal string to OauthToken struct failed, err: %v", err)
			res = &pb.GetTokenRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			return res, nil
		}
		return res, nil
	}
}
