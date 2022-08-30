package logic

import (
	"context"
	"github.com/tidwall/gjson"
	"main/app/common/log"
	"main/app/service/user/api/internal/svc"
	"main/app/service/user/api/internal/types"
	"main/app/service/user/rpc/crud/crud"
	"net/http"

	"github.com/imroc/req/v3"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(requ *types.LoginReq) (resp *types.LoginRes, err error) {
	logger := log.GetSugaredLogger()

	res, err := l.svcCtx.CrudRpcClient.Login(l.ctx, &crud.LoginReq{
		Username: requ.Username,
		Password: requ.Password,
		LastIp:   cast.ToString(l.ctx.Value("lastIp")),
	})
	if err != nil {
		logger.Errorf("login failed, err: %v", err)
		return &types.LoginRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}
	if res.Data == nil {
		return &types.LoginRes{
			Code: res.Code,
			Msg:  res.Msg,
			Ok:   res.Ok,
		}, nil
	}

	// 向 oauth 服务器请求签发 token
	resBody, err := req.NewRequest().SetHeader("Authorization", res.Data.AuthToken).
		Post("https://" + l.svcCtx.Domain + "/api/oauth/token/get")
	if err != nil {
		logger.Errorf("login failed, err: %v", err)
		return &types.LoginRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}, nil
	}
	if resBody.StatusCode != http.StatusOK {
		logger.Errorf("login failed, err: %v", err)
		return &types.LoginRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
			Ok:   false,
		}, nil
	}
	j := gjson.Parse(resBody.String())
	accessTokenValue := j.Get("data.access_token.token_value").String()
	refreshTokenValue := j.Get("data.access_token.refresh_token.token_value").String()

	return &types.LoginRes{
		Code: http.StatusOK,
		Msg:  "login successfully",
		Ok:   true,
		Data: types.LoginResData{
			AccessToken:  cast.ToString(accessTokenValue),
			RefreshToken: cast.ToString(refreshTokenValue),
		},
	}, nil
}
