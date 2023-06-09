// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

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

// LinkServiceClient is the client API for LinkService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LinkServiceClient interface {
	ShortenLink(ctx context.Context, in *LinkRequest, opts ...grpc.CallOption) (*LinkResponse, error)
	ExpandLink(ctx context.Context, in *LinkRequest, opts ...grpc.CallOption) (*LinkResponse, error)
}

type linkServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLinkServiceClient(cc grpc.ClientConnInterface) LinkServiceClient {
	return &linkServiceClient{cc}
}

func (c *linkServiceClient) ShortenLink(ctx context.Context, in *LinkRequest, opts ...grpc.CallOption) (*LinkResponse, error) {
	out := new(LinkResponse)
	err := c.cc.Invoke(ctx, "/linkgenerator.LinkService/ShortenLink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *linkServiceClient) ExpandLink(ctx context.Context, in *LinkRequest, opts ...grpc.CallOption) (*LinkResponse, error) {
	out := new(LinkResponse)
	err := c.cc.Invoke(ctx, "/linkgenerator.LinkService/ExpandLink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LinkServiceServer is the server API for LinkService service.
// All implementations must embed UnimplementedLinkServiceServer
// for forward compatibility
type LinkServiceServer interface {
	ShortenLink(context.Context, *LinkRequest) (*LinkResponse, error)
	ExpandLink(context.Context, *LinkRequest) (*LinkResponse, error)
	mustEmbedUnimplementedLinkServiceServer()
}

// UnimplementedLinkServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLinkServiceServer struct {
}

func (UnimplementedLinkServiceServer) ShortenLink(context.Context, *LinkRequest) (*LinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShortenLink not implemented")
}
func (UnimplementedLinkServiceServer) ExpandLink(context.Context, *LinkRequest) (*LinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExpandLink not implemented")
}
func (UnimplementedLinkServiceServer) mustEmbedUnimplementedLinkServiceServer() {}

// UnsafeLinkServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LinkServiceServer will
// result in compilation errors.
type UnsafeLinkServiceServer interface {
	mustEmbedUnimplementedLinkServiceServer()
}

func RegisterLinkServiceServer(s grpc.ServiceRegistrar, srv LinkServiceServer) {
	s.RegisterService(&LinkService_ServiceDesc, srv)
}

func _LinkService_ShortenLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LinkServiceServer).ShortenLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/linkgenerator.LinkService/ShortenLink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LinkServiceServer).ShortenLink(ctx, req.(*LinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LinkService_ExpandLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LinkServiceServer).ExpandLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/linkgenerator.LinkService/ExpandLink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LinkServiceServer).ExpandLink(ctx, req.(*LinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LinkService_ServiceDesc is the grpc.ServiceDesc for LinkService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LinkService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "linkgenerator.LinkService",
	HandlerType: (*LinkServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ShortenLink",
			Handler:    _LinkService_ShortenLink_Handler,
		},
		{
			MethodName: "ExpandLink",
			Handler:    _LinkService_ExpandLink_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/link_service.proto",
}
