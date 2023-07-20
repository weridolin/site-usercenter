// Code generated by goctl. DO NOT EDIT.
// Source: usercenter.proto

package usercenter

import (
	"context"

	"github.com/weridolin/site-gateway/services/users/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	GetUserInfoReq                = pb.GetUserInfoReq
	GetUserInfoResp               = pb.GetUserInfoResp
	GetUserMenuPermissionReq      = pb.GetUserMenuPermissionReq
	GetUserMenuPermissionResp     = pb.GetUserMenuPermissionResp
	GetUserResourcePermissionReq  = pb.GetUserResourcePermissionReq
	GetUserResourcePermissionResp = pb.GetUserResourcePermissionResp
	MenuPermissions               = pb.MenuPermissions
	ResourcePermissions           = pb.ResourcePermissions
	TokenValidateReq              = pb.TokenValidateReq
	TokenValidateResp             = pb.TokenValidateResp

	Usercenter interface {
		GetUserInfo(ctx context.Context, in *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoResp, error)
		GetUserResourcePermission(ctx context.Context, in *GetUserResourcePermissionReq, opts ...grpc.CallOption) (*GetUserResourcePermissionResp, error)
		GetUserMenuPermission(ctx context.Context, in *GetUserMenuPermissionReq, opts ...grpc.CallOption) (*GetUserMenuPermissionResp, error)
		TokenValidate(ctx context.Context, in *TokenValidateReq, opts ...grpc.CallOption) (*TokenValidateResp, error)
	}

	defaultUsercenter struct {
		cli zrpc.Client
	}
)

func NewUsercenter(cli zrpc.Client) Usercenter {
	return &defaultUsercenter{
		cli: cli,
	}
}

func (m *defaultUsercenter) GetUserInfo(ctx context.Context, in *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.GetUserInfo(ctx, in, opts...)
}

func (m *defaultUsercenter) GetUserResourcePermission(ctx context.Context, in *GetUserResourcePermissionReq, opts ...grpc.CallOption) (*GetUserResourcePermissionResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.GetUserResourcePermission(ctx, in, opts...)
}

func (m *defaultUsercenter) GetUserMenuPermission(ctx context.Context, in *GetUserMenuPermissionReq, opts ...grpc.CallOption) (*GetUserMenuPermissionResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.GetUserMenuPermission(ctx, in, opts...)
}

func (m *defaultUsercenter) TokenValidate(ctx context.Context, in *TokenValidateReq, opts ...grpc.CallOption) (*TokenValidateResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.TokenValidate(ctx, in, opts...)
}
