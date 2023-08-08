package user

import (
	"context"

	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/types"
	models "github.com/weridolin/site-gateway/services/users/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type TokenRefreshLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTokenRefreshLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TokenRefreshLogic {
	return &TokenRefreshLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TokenRefreshLogic) TokenRefresh(req *types.RefreshTokenReq) (resp *types.RefreshTokenResp, err error) {
	//expire exist token
	claims, err := models.ParseToken(req.Authorization, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		return nil, err
	}
	userID := int(claims["id"].(float64))
	user, err := l.svcCtx.UserModel.QueryUser(map[string]interface{}{"id": userID}, l.svcCtx.DB)
	if err != nil {
		// TODO return 404
		return nil, err
	}
	newAccessToken := models.GenToken(user, l.svcCtx.Config.JwtAuth.AccessSecret)

	return &types.RefreshTokenResp{
		BaseResponse: types.BaseResponse{
			Code: 0,
			Msg:  "刷新token成功",
		},
		Data: types.RefreshTokenPayload{
			AccessToken: newAccessToken,
		},
	}, nil
}
