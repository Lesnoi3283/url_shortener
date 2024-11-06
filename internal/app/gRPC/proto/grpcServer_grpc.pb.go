// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.6.1
// source: proto/grpcServer.proto

package proto

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	URLShortenerService_DeleteURLs_FullMethodName     = "/grpc_server.URLShortenerService/DeleteURLs"
	URLShortenerService_GetOriginalURL_FullMethodName = "/grpc_server.URLShortenerService/GetOriginalURL"
	URLShortenerService_PingDB_FullMethodName         = "/grpc_server.URLShortenerService/PingDB"
	URLShortenerService_Shorten_FullMethodName        = "/grpc_server.URLShortenerService/Shorten"
	URLShortenerService_ShortenBatch_FullMethodName   = "/grpc_server.URLShortenerService/ShortenBatch"
	URLShortenerService_Stats_FullMethodName          = "/grpc_server.URLShortenerService/Stats"
	URLShortenerService_UserURLs_FullMethodName       = "/grpc_server.URLShortenerService/UserURLs"
)

// URLShortenerServiceClient is the client API for URLShortenerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type URLShortenerServiceClient interface {
	DeleteURLs(ctx context.Context, in *DeleteURLsRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	GetOriginalURL(ctx context.Context, in *GetOriginalURLRequest, opts ...grpc.CallOption) (*GetAnOriginalURLResponse, error)
	PingDB(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error)
	Shorten(ctx context.Context, in *ShortenRequest, opts ...grpc.CallOption) (*ShortenResponse, error)
	ShortenBatch(ctx context.Context, in *ShortenBatchRequest, opts ...grpc.CallOption) (*ShortenBatchResponse, error)
	Stats(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*StatsResponse, error)
	UserURLs(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*UsersURLsResponse, error)
}

type uRLShortenerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewURLShortenerServiceClient(cc grpc.ClientConnInterface) URLShortenerServiceClient {
	return &uRLShortenerServiceClient{cc}
}

func (c *uRLShortenerServiceClient) DeleteURLs(ctx context.Context, in *DeleteURLsRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, URLShortenerService_DeleteURLs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceClient) GetOriginalURL(ctx context.Context, in *GetOriginalURLRequest, opts ...grpc.CallOption) (*GetAnOriginalURLResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAnOriginalURLResponse)
	err := c.cc.Invoke(ctx, URLShortenerService_GetOriginalURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceClient) PingDB(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, URLShortenerService_PingDB_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceClient) Shorten(ctx context.Context, in *ShortenRequest, opts ...grpc.CallOption) (*ShortenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShortenResponse)
	err := c.cc.Invoke(ctx, URLShortenerService_Shorten_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceClient) ShortenBatch(ctx context.Context, in *ShortenBatchRequest, opts ...grpc.CallOption) (*ShortenBatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShortenBatchResponse)
	err := c.cc.Invoke(ctx, URLShortenerService_ShortenBatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceClient) Stats(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*StatsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatsResponse)
	err := c.cc.Invoke(ctx, URLShortenerService_Stats_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceClient) UserURLs(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*UsersURLsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UsersURLsResponse)
	err := c.cc.Invoke(ctx, URLShortenerService_UserURLs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// URLShortenerServiceServer is the server API for URLShortenerService service.
// All implementations must embed UnimplementedURLShortenerServiceServer
// for forward compatibility.
type URLShortenerServiceServer interface {
	DeleteURLs(context.Context, *DeleteURLsRequest) (*empty.Empty, error)
	GetOriginalURL(context.Context, *GetOriginalURLRequest) (*GetAnOriginalURLResponse, error)
	PingDB(context.Context, *empty.Empty) (*empty.Empty, error)
	Shorten(context.Context, *ShortenRequest) (*ShortenResponse, error)
	ShortenBatch(context.Context, *ShortenBatchRequest) (*ShortenBatchResponse, error)
	Stats(context.Context, *empty.Empty) (*StatsResponse, error)
	UserURLs(context.Context, *empty.Empty) (*UsersURLsResponse, error)
	mustEmbedUnimplementedURLShortenerServiceServer()
}

// UnimplementedURLShortenerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedURLShortenerServiceServer struct{}

func (UnimplementedURLShortenerServiceServer) DeleteURLs(context.Context, *DeleteURLsRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteURLs not implemented")
}
func (UnimplementedURLShortenerServiceServer) GetOriginalURL(context.Context, *GetOriginalURLRequest) (*GetAnOriginalURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOriginalURL not implemented")
}
func (UnimplementedURLShortenerServiceServer) PingDB(context.Context, *empty.Empty) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingDB not implemented")
}
func (UnimplementedURLShortenerServiceServer) Shorten(context.Context, *ShortenRequest) (*ShortenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shorten not implemented")
}
func (UnimplementedURLShortenerServiceServer) ShortenBatch(context.Context, *ShortenBatchRequest) (*ShortenBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShortenBatch not implemented")
}
func (UnimplementedURLShortenerServiceServer) Stats(context.Context, *empty.Empty) (*StatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stats not implemented")
}
func (UnimplementedURLShortenerServiceServer) UserURLs(context.Context, *empty.Empty) (*UsersURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserURLs not implemented")
}
func (UnimplementedURLShortenerServiceServer) mustEmbedUnimplementedURLShortenerServiceServer() {}
func (UnimplementedURLShortenerServiceServer) testEmbeddedByValue()                             {}

// UnsafeURLShortenerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to URLShortenerServiceServer will
// result in compilation errors.
type UnsafeURLShortenerServiceServer interface {
	mustEmbedUnimplementedURLShortenerServiceServer()
}

func RegisterURLShortenerServiceServer(s grpc.ServiceRegistrar, srv URLShortenerServiceServer) {
	// If the following call pancis, it indicates UnimplementedURLShortenerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&URLShortenerService_ServiceDesc, srv)
}

func _URLShortenerService_DeleteURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteURLsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).DeleteURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_DeleteURLs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).DeleteURLs(ctx, req.(*DeleteURLsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerService_GetOriginalURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOriginalURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).GetOriginalURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_GetOriginalURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).GetOriginalURL(ctx, req.(*GetOriginalURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerService_PingDB_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).PingDB(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_PingDB_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).PingDB(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerService_Shorten_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).Shorten(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_Shorten_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).Shorten(ctx, req.(*ShortenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerService_ShortenBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortenBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).ShortenBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_ShortenBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).ShortenBatch(ctx, req.(*ShortenBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerService_Stats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).Stats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_Stats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).Stats(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerService_UserURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).UserURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_UserURLs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).UserURLs(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// URLShortenerService_ServiceDesc is the grpc.ServiceDesc for URLShortenerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var URLShortenerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc_server.URLShortenerService",
	HandlerType: (*URLShortenerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteURLs",
			Handler:    _URLShortenerService_DeleteURLs_Handler,
		},
		{
			MethodName: "GetOriginalURL",
			Handler:    _URLShortenerService_GetOriginalURL_Handler,
		},
		{
			MethodName: "PingDB",
			Handler:    _URLShortenerService_PingDB_Handler,
		},
		{
			MethodName: "Shorten",
			Handler:    _URLShortenerService_Shorten_Handler,
		},
		{
			MethodName: "ShortenBatch",
			Handler:    _URLShortenerService_ShortenBatch_Handler,
		},
		{
			MethodName: "Stats",
			Handler:    _URLShortenerService_Stats_Handler,
		},
		{
			MethodName: "UserURLs",
			Handler:    _URLShortenerService_UserURLs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/grpcServer.proto",
}
