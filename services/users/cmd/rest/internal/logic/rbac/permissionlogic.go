package rbac

import (
	"context"
	"fmt"

	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/types"
	"github.com/weridolin/site-gateway/services/users/models"
	"github.com/weridolin/site-gateway/tools"
	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionLogic {
	return &PermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PermissionLogic) ReLoadResourceLogic(req *types.ReLoadResourceReq) (resp *types.ReLoadResourceResp, err error) {
	// 先清除缓存
	l.svcCtx.RedisClient.Del(l.ctx, tools.ResourceAuthenticatedCacheKey)
	// 获取所有的用户id
	var user_list []models.User
	err = l.svcCtx.DB.Find(&user_list).Error
	if err != nil {
		return nil, err
	}
	// 清楚所有用户的缓存
	for _, user := range user_list {
		cache_key := tools.UserPermissionKey(user.ID)
		l.svcCtx.RedisClient.Del(l.ctx, cache_key)
		fmt.Println("清除用户缓存 -> ", cache_key)
	}

	// 重新加载资源
	svc.LoadInitData(l.svcCtx.Config, l.svcCtx.DB)

	return &types.ReLoadResourceResp{
		BaseResponse: types.BaseResponse{
			Code: 200,
			Msg:  "重新加载资源权限成功",
		},
	}, nil
}
