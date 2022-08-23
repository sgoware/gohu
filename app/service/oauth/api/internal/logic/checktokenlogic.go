package logic

import (
	"context"
	"main/app/service/oauth/api/internal/svc"
	"main/app/service/oauth/api/internal/token"
	"main/app/service/oauth/api/internal/types"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckTokenLogic {
	return &CheckTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckTokenLogic) CheckToken(req *types.CheckTokenReq) (*types.CheckTokenRes, error) {
	tokenService := token.GetTokenService()
	oauthToken, err := tokenService.ReadAccessToken(l.ctx, req.OAtuh2Token)
	if err != nil {
		logx.Errorf("parse token failed, err: %v", err)
		return &types.CheckTokenRes{
			Code: http.StatusBadRequest,
			Msg:  "invalid token",
			Ok:   false,
		}, nil
	}
	if oauthToken.TokenType != req.TokenType {
		return &types.CheckTokenRes{
			Code: http.StatusOK,
			Msg:  "incorrect token type",
			Ok:   false,
		}, nil
	}
	// TODO: 待校验
	if oauthToken.ExpiresAt < time.Now().Unix() {
		return &types.CheckTokenRes{
			Code: http.StatusBadRequest,
			Msg:  "token is expired",
			Ok:   false,
		}, nil
	}
	return &types.CheckTokenRes{
		Code: http.StatusOK,
		Msg:  "token is valid",
		Ok:   true,
	}, nil
}
