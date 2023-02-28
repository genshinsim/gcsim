// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: protos/backend/queue.proto

package queue

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

// WorkQueueClient is the client API for WorkQueue service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WorkQueueClient interface {
	Get(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*GetResp, error)
	Complete(ctx context.Context, in *CompleteReq, opts ...grpc.CallOption) (*CompleteResp, error)
}

type workQueueClient struct {
	cc grpc.ClientConnInterface
}

func NewWorkQueueClient(cc grpc.ClientConnInterface) WorkQueueClient {
	return &workQueueClient{cc}
}

func (c *workQueueClient) Get(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*GetResp, error) {
	out := new(GetResp)
	err := c.cc.Invoke(ctx, "/queue.WorkQueue/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workQueueClient) Complete(ctx context.Context, in *CompleteReq, opts ...grpc.CallOption) (*CompleteResp, error) {
	out := new(CompleteResp)
	err := c.cc.Invoke(ctx, "/queue.WorkQueue/Complete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WorkQueueServer is the server API for WorkQueue service.
// All implementations must embed UnimplementedWorkQueueServer
// for forward compatibility
type WorkQueueServer interface {
	Get(context.Context, *GetReq) (*GetResp, error)
	Complete(context.Context, *CompleteReq) (*CompleteResp, error)
	mustEmbedUnimplementedWorkQueueServer()
}

// UnimplementedWorkQueueServer must be embedded to have forward compatible implementations.
type UnimplementedWorkQueueServer struct {
}

func (UnimplementedWorkQueueServer) Get(context.Context, *GetReq) (*GetResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedWorkQueueServer) Complete(context.Context, *CompleteReq) (*CompleteResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Complete not implemented")
}
func (UnimplementedWorkQueueServer) mustEmbedUnimplementedWorkQueueServer() {}

// UnsafeWorkQueueServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WorkQueueServer will
// result in compilation errors.
type UnsafeWorkQueueServer interface {
	mustEmbedUnimplementedWorkQueueServer()
}

func RegisterWorkQueueServer(s grpc.ServiceRegistrar, srv WorkQueueServer) {
	s.RegisterService(&WorkQueue_ServiceDesc, srv)
}

func _WorkQueue_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkQueueServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/queue.WorkQueue/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkQueueServer).Get(ctx, req.(*GetReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkQueue_Complete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompleteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkQueueServer).Complete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/queue.WorkQueue/Complete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkQueueServer).Complete(ctx, req.(*CompleteReq))
	}
	return interceptor(ctx, in, info, handler)
}

// WorkQueue_ServiceDesc is the grpc.ServiceDesc for WorkQueue service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WorkQueue_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "queue.WorkQueue",
	HandlerType: (*WorkQueueServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _WorkQueue_Get_Handler,
		},
		{
			MethodName: "Complete",
			Handler:    _WorkQueue_Complete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/backend/queue.proto",
}