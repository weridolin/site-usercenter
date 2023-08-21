package user

import (
	"context"
	"fmt"

	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/types"
	"github.com/weridolin/site-gateway/services/users/models"
	"github.com/zeromicro/go-zero/core/logx"
)

type RbacLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRbacLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RbacLogic {
	return &RbacLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RbacLogic) GetMenu(userId int, req *types.GetMenuReq) (resp *types.GetMenuResp, err error) {
	var roleIds []int
	roleRefs, _ := models.UserRoles{}.Query(map[string]interface{}{"user_id": userId}, l.svcCtx.DB)
	for _, v := range roleRefs {
		roleIds = append(roleIds, v.RoleId)
	}
	fmt.Println("roleIds --> ", roleIds)
	if roleIds != nil {
		roles, err := models.QueryRolePermission(roleIds, l.svcCtx.DB)
		if err != nil {
			fmt.Println("get menu error", err)
			return nil, err
		}
		var menus []*models.Menu
		for _, v := range roles {
			fmt.Println("role --> ", v.Menus)
			menus = append(menus, v.Menus...)
		}
		return &types.GetMenuResp{
			BaseResponse: types.BaseResponse{
				Code: 0,
				Msg:  "获取成功",
			},
			Data: types.Menu{}.FromMenuModels(menus),
		}, nil
	}
	return &types.GetMenuResp{
		BaseResponse: types.BaseResponse{
			Code: 0,
			Msg:  "获取成功",
		},
		Data: nil,
	}, nil
}
