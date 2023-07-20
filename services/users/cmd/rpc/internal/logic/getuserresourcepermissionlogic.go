package logic

import (
	"context"

	"github.com/weridolin/site-gateway/services/users/cmd/rpc/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rpc/pb"
	"github.com/weridolin/site-gateway/services/users/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserResourcePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserResourcePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserResourcePermissionLogic {
	return &GetUserResourcePermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserResourcePermissionLogic) GetUserResourcePermission(in *pb.GetUserResourcePermissionReq) (*pb.GetUserResourcePermissionResp, error) {
	resources_res, err := models.GetUserResourcePermissionByUserID(l.svcCtx.DB, int(in.UserId))
	if err != nil {
		return nil, err
	} else {
		var resources []*pb.ResourcePermissions
		for _, _resources := range resources_res {
			perm := &pb.ResourcePermissions{
				Id:          int64(_resources.ID),
				ServerName:  _resources.ServerName,
				Version:     _resources.Version,
				Method:      _resources.Method,
				Description: _resources.Description,
				Url:         _resources.Url,
			}
			resources = append(resources, perm)
		}
		return &pb.GetUserResourcePermissionResp{
			ResourcePermissions: resources,
		}, nil
	}

	// return &pb.GetUserResourcePermissionResp{}, nil
}
