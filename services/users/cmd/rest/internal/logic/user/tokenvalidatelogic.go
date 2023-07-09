package user

import (
	"context"
	"fmt"

	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/types"
	"github.com/weridolin/site-gateway/services/users/models"
	"github.com/weridolin/site-gateway/tools"
	"github.com/zeromicro/go-zero/core/logx"
	xerrors "github.com/zeromicro/x/errors"
)

type TokenValidateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTokenValidateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TokenValidateLogic {
	return &TokenValidateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TokenValidateLogic) TokenValidate(req *types.ValidateTokenReq) (resp *types.ValidateResp, err error) {
	// userID := l.ctx.Value("id")
	err = l.TokenUnregister(req.Authorization)
	if err != nil {
		return nil, err
	}
	// 校验token是否合法
	claims, err := models.ParseToken(req.Authorization, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		return nil, err
	}
	userID := claims["id"]
	return &types.ValidateResp{
		UserId: fmt.Sprintf("%v", userID.(float64)),
	}, nil
}

func (l *TokenValidateLogic) TokenUnregister(token string) error {
	key := tools.InvalidTokenKey(token)
	res := l.svcCtx.RedisClient.Get(key)
	if res.Err() != nil {
		// 不存在，说明token没有被注销
		return nil
	} else {
		return xerrors.New(401, "token is invalid")
	}
}
