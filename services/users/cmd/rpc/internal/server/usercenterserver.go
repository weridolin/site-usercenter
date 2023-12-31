// Code generated by goctl. DO NOT EDIT.
// Source: usercenter.proto

package server

import (
	"context"

	"github.com/weridolin/site-gateway/services/users/cmd/rpc/internal/logic"
	"github.com/weridolin/site-gateway/services/users/cmd/rpc/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rpc/pb"
)

type UsercenterServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedUsercenterServer
}

func NewUsercenterServer(svcCtx *svc.ServiceContext) *UsercenterServer {
	return &UsercenterServer{
		svcCtx: svcCtx,
	}
}

func (s *UsercenterServer) GetUserInfo(ctx context.Context, in *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	l := logic.NewGetUserInfoLogic(ctx, s.svcCtx)
	return l.GetUserInfo(in)
}

func (s *UsercenterServer) GetUserResourcePermission(ctx context.Context, in *pb.GetUserResourcePermissionReq) (*pb.GetUserResourcePermissionResp, error) {
	l := logic.NewGetUserResourcePermissionLogic(ctx, s.svcCtx)
	return l.GetUserResourcePermission(in)
}

func (s *UsercenterServer) GetUserMenuPermission(ctx context.Context, in *pb.GetUserMenuPermissionReq) (*pb.GetUserMenuPermissionResp, error) {
	l := logic.NewGetUserMenuPermissionLogic(ctx, s.svcCtx)
	return l.GetUserMenuPermission(in)
}

func (s *UsercenterServer) TokenValidate(ctx context.Context, in *pb.TokenValidateReq) (*pb.TokenValidateResp, error) {
	l := logic.NewTokenValidateLogic(ctx, s.svcCtx)
	return l.TokenValidate(in)
}

func (s *UsercenterServer) GetMutipleUserInfo(ctx context.Context, in *pb.GetMutipleUserInfoReq) (*pb.GetMutipleUserInfoResp, error) {
	l := logic.NewGetMutipleUserInfoLogic(ctx, s.svcCtx)
	return l.GetMutipleUserInfo(in)
}
