package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
)

func NewIPInterceptor(allowedNet *net.IPNet) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if info.FullMethod == "/grpc_server.URLShortenerService/Stats" {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
			}

			values := md.Get("X-Real-IP")
			if len(values) == 0 {
				return nil, status.Errorf(codes.Unauthenticated, "missing IP address")
			}

			ip := net.ParseIP(values[0])
			if ip == nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid IP address")
			}

			if !allowedNet.Contains(ip) {
				return nil, status.Errorf(codes.PermissionDenied, "access denied")
			}
		}

		return handler(ctx, req)
	}
}
