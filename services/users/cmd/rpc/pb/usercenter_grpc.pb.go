// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.19.4
// source: pb/usercenter.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Usercenter_GetUserInfo_FullMethodName               = "/pb.usercenter/getUserInfo"
	Usercenter_GetUserResourcePermission_FullMethodName = "/pb.usercenter/getUserResourcePermission"
	Usercenter_GetUserMenuPermission_FullMethodName     = "/pb.usercenter/getUserMenuPermission"
	Usercenter_TokenValidate_FullMethodName             = "/pb.usercenter/tokenValidate"
	Usercenter_GetMutipleUserInfo_FullMethodName        = "/pb.usercenter/getMutipleUserInfo"
)

// UsercenterClient is the client API for Usercenter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UsercenterClient interface {
	GetUserInfo(ctx context.Context, in *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoResp, error)
	GetUserResourcePermission(ctx context.Context, in *GetUserResourcePermissionReq, opts ...grpc.CallOption) (*GetUserResourcePermissionResp, error)
	GetUserMenuPermission(ctx context.Context, in *GetUserMenuPermissionReq, opts ...grpc.CallOption) (*GetUserMenuPermissionResp, error)
	TokenValidate(ctx context.Context, in *TokenValidateReq, opts ...grpc.CallOption) (*TokenValidateResp, error)
	GetMutipleUserInfo(ctx context.Context, in *GetMutipleUserInfoReq, opts ...grpc.CallOption) (*GetMutipleUserInfoResp, error)
}

type usercenterClient struct {
	cc grpc.ClientConnInterface
}

func NewUsercenterClient(cc grpc.ClientConnInterface) UsercenterClient {
	return &usercenterClient{cc}
}

func (c *usercenterClient) GetUserInfo(ctx context.Context, in *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoResp, error) {
	out := new(GetUserInfoResp)
	err := c.cc.Invoke(ctx, Usercenter_GetUserInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usercenterClient) GetUserResourcePermission(ctx context.Context, in *GetUserResourcePermissionReq, opts ...grpc.CallOption) (*GetUserResourcePermissionResp, error) {
	out := new(GetUserResourcePermissionResp)
	err := c.cc.Invoke(ctx, Usercenter_GetUserResourcePermission_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usercenterClient) GetUserMenuPermission(ctx context.Context, in *GetUserMenuPermissionReq, opts ...grpc.CallOption) (*GetUserMenuPermissionResp, error) {
	out := new(GetUserMenuPermissionResp)
	err := c.cc.Invoke(ctx, Usercenter_GetUserMenuPermission_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usercenterClient) TokenValidate(ctx context.Context, in *TokenValidateReq, opts ...grpc.CallOption) (*TokenValidateResp, error) {
	out := new(TokenValidateResp)
	err := c.cc.Invoke(ctx, Usercenter_TokenValidate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usercenterClient) GetMutipleUserInfo(ctx context.Context, in *GetMutipleUserInfoReq, opts ...grpc.CallOption) (*GetMutipleUserInfoResp, error) {
	out := new(GetMutipleUserInfoResp)
	err := c.cc.Invoke(ctx, Usercenter_GetMutipleUserInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsercenterServer is the server API for Usercenter service.
// All implementations must embed UnimplementedUsercenterServer
// for forward compatibility
type UsercenterServer interface {
	GetUserInfo(context.Context, *GetUserInfoReq) (*GetUserInfoResp, error)
	GetUserResourcePermission(context.Context, *GetUserResourcePermissionReq) (*GetUserResourcePermissionResp, error)
	GetUserMenuPermission(context.Context, *GetUserMenuPermissionReq) (*GetUserMenuPermissionResp, error)
	TokenValidate(context.Context, *TokenValidateReq) (*TokenValidateResp, error)
	GetMutipleUserInfo(context.Context, *GetMutipleUserInfoReq) (*GetMutipleUserInfoResp, error)
	mustEmbedUnimplementedUsercenterServer()
}

// UnimplementedUsercenterServer must be embedded to have forward compatible implementations.
type UnimplementedUsercenterServer struct {
}

func (UnimplementedUsercenterServer) GetUserInfo(context.Context, *GetUserInfoReq) (*GetUserInfoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfo not implemented")
}
func (UnimplementedUsercenterServer) GetUserResourcePermission(context.Context, *GetUserResourcePermissionReq) (*GetUserResourcePermissionResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserResourcePermission not implemented")
}
func (UnimplementedUsercenterServer) GetUserMenuPermission(context.Context, *GetUserMenuPermissionReq) (*GetUserMenuPermissionResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserMenuPermission not implemented")
}
func (UnimplementedUsercenterServer) TokenValidate(context.Context, *TokenValidateReq) (*TokenValidateResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TokenValidate not implemented")
}
func (UnimplementedUsercenterServer) GetMutipleUserInfo(context.Context, *GetMutipleUserInfoReq) (*GetMutipleUserInfoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMutipleUserInfo not implemented")
}
func (UnimplementedUsercenterServer) mustEmbedUnimplementedUsercenterServer() {}

// UnsafeUsercenterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UsercenterServer will
// result in compilation errors.
type UnsafeUsercenterServer interface {
	mustEmbedUnimplementedUsercenterServer()
}

func RegisterUsercenterServer(s grpc.ServiceRegistrar, srv UsercenterServer) {
	s.RegisterService(&Usercenter_ServiceDesc, srv)
}

func _Usercenter_GetUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsercenterServer).GetUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Usercenter_GetUserInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsercenterServer).GetUserInfo(ctx, req.(*GetUserInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Usercenter_GetUserResourcePermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserResourcePermissionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsercenterServer).GetUserResourcePermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Usercenter_GetUserResourcePermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsercenterServer).GetUserResourcePermission(ctx, req.(*GetUserResourcePermissionReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Usercenter_GetUserMenuPermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserMenuPermissionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsercenterServer).GetUserMenuPermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Usercenter_GetUserMenuPermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsercenterServer).GetUserMenuPermission(ctx, req.(*GetUserMenuPermissionReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Usercenter_TokenValidate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenValidateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsercenterServer).TokenValidate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Usercenter_TokenValidate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsercenterServer).TokenValidate(ctx, req.(*TokenValidateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Usercenter_GetMutipleUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMutipleUserInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsercenterServer).GetMutipleUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Usercenter_GetMutipleUserInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsercenterServer).GetMutipleUserInfo(ctx, req.(*GetMutipleUserInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Usercenter_ServiceDesc is the grpc.ServiceDesc for Usercenter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Usercenter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.usercenter",
	HandlerType: (*UsercenterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getUserInfo",
			Handler:    _Usercenter_GetUserInfo_Handler,
		},
		{
			MethodName: "getUserResourcePermission",
			Handler:    _Usercenter_GetUserResourcePermission_Handler,
		},
		{
			MethodName: "getUserMenuPermission",
			Handler:    _Usercenter_GetUserMenuPermission_Handler,
		},
		{
			MethodName: "tokenValidate",
			Handler:    _Usercenter_TokenValidate_Handler,
		},
		{
			MethodName: "getMutipleUserInfo",
			Handler:    _Usercenter_GetMutipleUserInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/usercenter.proto",
}
