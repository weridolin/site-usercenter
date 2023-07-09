package user

import (
	"context"
	"time"

	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/types"
	"github.com/weridolin/site-gateway/services/users/models"
	"github.com/weridolin/site-gateway/tools"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.ValidateTokenReq) (resp *types.LogoutResp, err error) {
	// 校验token是否合法
	claims, err := models.ParseToken(req.Authorization, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		return nil, err
	}
	timeRemain := time.Unix(int64(claims["exp"].(float64)), 0).Sub(time.Now())
	l.svcCtx.RedisClient.Set(tools.InvalidTokenKey(req.Authorization), 1, timeRemain*time.Second)
	return &types.LogoutResp{
		BaseResponse: types.BaseResponse{
			Code: 0,
			Msg:  "登出成功",
		},
	}, nil
}
