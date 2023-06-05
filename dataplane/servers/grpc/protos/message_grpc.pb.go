// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: message.proto

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

// MessageClient is the client API for Message service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageClient interface {
	Create(ctx context.Context, in *MessageCreateReq, opts ...grpc.CallOption) (*MessageCreateRes, error)
}

type messageClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageClient(cc grpc.ClientConnInterface) MessageClient {
	return &messageClient{cc}
}

func (c *messageClient) Create(ctx context.Context, in *MessageCreateReq, opts ...grpc.CallOption) (*MessageCreateRes, error) {
	out := new(MessageCreateRes)
	err := c.cc.Invoke(ctx, "/kanthor.dataplane.v1.Message/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessageServer is the server API for Message service.
// All implementations must embed UnimplementedMessageServer
// for forward compatibility
type MessageServer interface {
	Create(context.Context, *MessageCreateReq) (*MessageCreateRes, error)
	mustEmbedUnimplementedMessageServer()
}

// UnimplementedMessageServer must be embedded to have forward compatible implementations.
type UnimplementedMessageServer struct {
}

func (UnimplementedMessageServer) Create(context.Context, *MessageCreateReq) (*MessageCreateRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedMessageServer) mustEmbedUnimplementedMessageServer() {}

// UnsafeMessageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServer will
// result in compilation errors.
type UnsafeMessageServer interface {
	mustEmbedUnimplementedMessageServer()
}

func RegisterMessageServer(s grpc.ServiceRegistrar, srv MessageServer) {
	s.RegisterService(&Message_ServiceDesc, srv)
}

func _Message_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageCreateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kanthor.dataplane.v1.Message/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServer).Create(ctx, req.(*MessageCreateReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Message_ServiceDesc is the grpc.ServiceDesc for Message service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Message_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kanthor.dataplane.v1.Message",
	HandlerType: (*MessageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Message_Create_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "message.proto",
}
