package logic

import (
	"context"

	"github.com/weridolin/site-gateway/services/users/cmd/rpc/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rpc/pb"
	"github.com/weridolin/site-gateway/services/users/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserMenuPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserMenuPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserMenuPermissionLogic {
	return &GetUserMenuPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserMenuPermissionLogic) GetUserMenuPermission(in *pb.GetUserMenuPermissionReq) (*pb.GetUserMenuPermissionResp, error) {
	// user,err:=l.svcCtx.UserModel.QueryUser(map[string]interface{}{"id":in.UserId},l.svcCtx.DB)
	menus_res, err := models.GetUserMenuPermissionByUserID(l.svcCtx.DB, int(in.UserId))
	if err != nil {
		return nil, err
	} else {
		var menus []*pb.MenuPermissions
		for _, _menu := range menus_res {
			menu := &pb.MenuPermissions{
				Id:            int64(_menu.ID),
				MenuName:      _menu.Name,
				MenuComponent: _menu.Component,
				MenuIcon:      _menu.Icon,
				MenuUrl:       _menu.Url,
				MenuType:      int32(_menu.Type),
				MenuRouteName: _menu.RouteName,
				ParentId:      int64(_menu.ParentId),
			}
			menus = append(menus, menu)
		}

		return &pb.GetUserMenuPermissionResp{
			MenuPermissions: menus,
		}, nil
	}

}
