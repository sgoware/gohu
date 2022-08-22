package logic

import (
	"context"
	"encoding/json"
	"gohu/app/common/log"
	"gohu/app/service/oauth/model"
	"gohu/app/service/oauth/rpc/token/store/internal/svc"
	"gohu/app/service/oauth/rpc/token/store/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
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

func (l *GetTokenLogic) GetToken(in *pb.GetTokenReq) (*pb.GetTokenRes, error) {
	logger := log.GetSugaredLogger()

	if in.UserId == "" {
		logger.Errorf("get token failed, err: %v", model.ErrInvalidTokenRequest)
		return &pb.GetTokenRes{
			Ok:  false,
			Msg: "get token failed, err: invalid token request",
		}, nil
	}

	val, err := l.svcCtx.Rdb.Get(l.ctx, model.JwtToken+"_"+in.UserId).Result()
	if err == nil {
		res := &pb.GetTokenRes{
			Ok:   true,
			Msg:  "get token successfully",
			Data: &pb.GetTokenRes_Data{OauthToken: &pb.OAuth2Token{}},
		}
		err = json.Unmarshal([]byte(val), &res.Data.OauthToken)
		if err != nil {
			logger.Errorf("unmarshal string to OauthToken struct failed, err: %v", err)
			return &pb.GetTokenRes{
				Ok:   false,
				Msg:  "internal server err",
				Data: nil,
			}, nil
		}
		return res, nil
	} else {
		// TODO: 待解决redis err
		if err != redis.ErrEmptyKey {
			return &pb.GetTokenRes{
				Ok:  false,
				Msg: "get token failed, err: redis err",
			}, nil
		} else {
			return &pb.GetTokenRes{
				Ok:  false,
				Msg: "token not found",
			}, nil
		}
	}
}
