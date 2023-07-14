package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
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

func (l *TokenValidateLogic) TokenValidate(req *types.ValidateTokenReq, apiPermissionRequired string) (resp *types.ValidateResp, err error) {
	// userID := l.ctx.Value("id")
	err = l.TokenUnregister(req.Authorization)
	if err != nil {
		return nil, err
	}
	// 校验token是否合法
	claims, err := models.ParseToken(req.Authorization, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		fmt.Println("parse token error", err)
		return nil, xerrors.New(http.StatusUnauthorized, err.Error())
	}
	userID := claims["id"]
	// 校验用户是否有权限
	if !l.HasPermission(int(userID.(float64)), apiPermissionRequired) {
		return nil, xerrors.New(http.StatusForbidden, "user has no permission")
	}
	return &types.ValidateResp{
		UserId: fmt.Sprintf("%v", userID.(float64)),
	}, nil
}

func (l *TokenValidateLogic) TokenUnregister(token string) error {
	key := tools.InvalidTokenKey(token)
	res := l.svcCtx.RedisClient.Get(l.ctx, key)
	if res.Err() != nil {
		// 不存在，说明token没有被注销
		return nil
	} else {
		return xerrors.New(http.StatusUnauthorized, "token is invalid")
	}
}

func (l *TokenValidateLogic) HasPermission(userId int, apiPermissionRequired string) bool {
	fmt.Println("api required permission", apiPermissionRequired, userId)
	key := tools.UserPermissionKey(userId)
	val, err := l.svcCtx.RedisClient.Get(l.ctx, key).Result()
	res := false
	if err != nil {
		if err == redis.Nil {
			resourceList := make([]string, 0)
			// 从数据库查询对应的权限
			roles, err := l.svcCtx.RoleModel.Query(map[string]interface{}{"create_user": userId, "zone": "site"}, l.svcCtx.DB)
			if err != nil {
				fmt.Println("get permission from db error", err)
			} else {
				for _, role := range roles {
					for _, resource := range role.Resources {
						if apiPermissionRequired == resource.ServerName+":"+resource.Url+":"+resource.Method {
							res = true
						}
						resourceList = append(resourceList, resource.ServerName+":"+resource.Url+":"+resource.Method)
					}
				}
				// 将权限存入redis
				resourceListJson, _ := json.Marshal(resourceList)
				res := l.svcCtx.RedisClient.Set(l.ctx, tools.UserPermissionKey(userId), resourceListJson, 0)
				fmt.Println("set permission to redis res", res)

			}
		} else {
			fmt.Println("get permission from cache err", err)
		}
	} else {
		var permission []string
		json.Unmarshal([]byte(val), &permission)
		for _, p := range permission {
			if p == apiPermissionRequired {
				res = true
				break
			}
		}
		fmt.Println("get permission from cache", permission)
	}
	return res
}
