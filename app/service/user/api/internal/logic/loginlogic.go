package logic

import (
	"context"
	"github.com/tidwall/gjson"
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
	res, err := l.svcCtx.CrudRpcClient.Login(l.ctx, &crud.LoginReq{
		Username: requ.Username,
		Password: requ.Password,
		LastIp:   cast.ToString(l.ctx.Value("lastIp")),
	})
	if err != nil {
		logx.Errorf("login failed, err: %v", err)
	}
	if res.Data == nil {
		return &types.LoginRes{
			Code: int(res.Code),
			Msg:  res.Msg,
		}, nil
	}

	// 向 oauth 服务器请求签发 token
	resBody, err := req.NewRequest().SetHeader("Authorization", res.Data.AuthToken).
		Post("https://" + l.svcCtx.Domain + "/api/oauth/token/get")
	if err != nil {
		logx.Infof("%v", err)
		return &types.LoginRes{
			Code: http.StatusInternalServerError,
			Msg:  "login failed, err: internal server err",
		}, nil
	}
	if resBody.StatusCode != http.StatusOK {
		logx.Infof("%v", res.String())
		return &types.LoginRes{
			Code: http.StatusInternalServerError,
			Msg:  "login failed, err: internal server err",
		}, nil
	}
	j := gjson.Parse(resBody.String())
	accessTokenValue := j.Get("data.access_token.token_value").String()
	refreshTokenValue := j.Get("data.access_token.refresh_token.token_value").String()

	return &types.LoginRes{
		Code: int(res.Code),
		Msg:  "login successfully",
		Data: types.LoginResData{
			AccessToken:  cast.ToString(accessTokenValue),
			RefreshToken: cast.ToString(refreshTokenValue),
		},
	}, nil
}
