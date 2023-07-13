// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: workspace.proto

package protos

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

// WsClient is the client API for Ws service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WsClient interface {
	Get(ctx context.Context, in *WsGetReq, opts ...grpc.CallOption) (*Workspace, error)
}

type wsClient struct {
	cc grpc.ClientConnInterface
}

func NewWsClient(cc grpc.ClientConnInterface) WsClient {
	return &wsClient{cc}
}

func (c *wsClient) Get(ctx context.Context, in *WsGetReq, opts ...grpc.CallOption) (*Workspace, error) {
	out := new(Workspace)
	err := c.cc.Invoke(ctx, "/kanthor.controlplane.v1.Ws/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WsServer is the server API for Ws service.
// All implementations must embed UnimplementedWsServer
// for forward compatibility
type WsServer interface {
	Get(context.Context, *WsGetReq) (*Workspace, error)
	mustEmbedUnimplementedWsServer()
}

// UnimplementedWsServer must be embedded to have forward compatible implementations.
type UnimplementedWsServer struct {
}

func (UnimplementedWsServer) Get(context.Context, *WsGetReq) (*Workspace, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedWsServer) mustEmbedUnimplementedWsServer() {}

// UnsafeWsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WsServer will
// result in compilation errors.
type UnsafeWsServer interface {
	mustEmbedUnimplementedWsServer()
}

func RegisterWsServer(s grpc.ServiceRegistrar, srv WsServer) {
	s.RegisterService(&Ws_ServiceDesc, srv)
}

func _Ws_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WsGetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WsServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kanthor.controlplane.v1.Ws/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WsServer).Get(ctx, req.(*WsGetReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Ws_ServiceDesc is the grpc.ServiceDesc for Ws service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ws_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kanthor.controlplane.v1.Ws",
	HandlerType: (*WsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Ws_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "workspace.proto",
}
