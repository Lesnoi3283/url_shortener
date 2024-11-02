package interceptors

import (
	"context"
	"github.com/Lesnoi3283/url_shortener/pkg/secure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	UserIDContextKey contextKey = "userID"
	NoUserIDValue    int        = -1
)

func NewUnaryAuthInterceptor(jh *secure.JWTHelper) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values := md.Get("token")
			if len(values) > 0 {
				token := values[0]
				userID, err := jh.GetUserID(token)
				if err != nil {
					return nil, status.Errorf(codes.Unauthenticated, "invalid token")
				}
				return handler(context.WithValue(ctx, UserIDContextKey, userID), req)
			}
		}

		return handler(context.WithValue(ctx, UserIDContextKey, -1), req)
	}
}
