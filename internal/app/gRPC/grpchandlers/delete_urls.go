package grpchandlers

import (
	"context"
	"github.com/Lesnoi3283/url_shortener/internal/app/gRPC/interceptors"
	"github.com/Lesnoi3283/url_shortener/internal/app/gRPC/proto"
	"github.com/Lesnoi3283/url_shortener/internal/app/logic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ShortenerServer) DeleteURLs(ctx context.Context, req *proto.DeleteURLsRequest) (*emptypb.Empty, error) {
	//auth
	userID := ctx.Value(interceptors.UserIDContextKey)
	if userID == nil || userID == -1 {
		s.Logger.Debug("UserID not found req ctx")
		return &emptypb.Empty{}, status.Errorf(codes.Unauthenticated, "User ID not found")
	}
	userIDInt, ok := userID.(int)
	if !ok {
		s.Logger.Warnf("User ID is not an int, real type: `%T`, value: `%v`", userID, userID)
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "User ID is not an int")
	}

	//delete urls
	err := logic.DeleteURLs(userIDInt, req.URLs, s.Storage)
	if err != nil {
		s.Logger.Errorf("DeleteURLs error: %v", err)
		return &emptypb.Empty{}, status.Error(codes.Internal, "Internal server error")
	}
	return &emptypb.Empty{}, nil
}
