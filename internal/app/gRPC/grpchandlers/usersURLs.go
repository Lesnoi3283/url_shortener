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

func (s *ShortenerServer) UserURLs(ctx context.Context, req *emptypb.Empty) (*proto.UsersURLsResponse, error) {
	//auth
	userID := ctx.Value(interceptors.UserIDContextKey)
	if userID == nil || userID == -1 {
		s.Logger.Debug("UserID not found req ctx")
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found")
	}
	userIDInt, ok := userID.(int)
	if !ok {
		s.Logger.Warnf("User ID is not an int, real type: `%T`, value: `%v`", userID, userID)
		return nil, status.Error(codes.InvalidArgument, "User ID is not an int")
	}

	//get user`s URLs
	URLs, err := logic.GetUsersURLs(ctx, s.Storage, s.Conf.BaseAddress, userIDInt)
	if err != nil {
		s.Logger.Errorf("GetUserURLs error: %v", err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	//prepare response
	response := &proto.UsersURLsResponse{
		Urls: make([]*proto.UsersURLsResponse_URL, len(URLs)),
	}
	for i, u := range URLs {
		response.Urls[i].Original = u.OriginalURL
		response.Urls[i].Short = u.ShortURL
	}

	//return response
	return response, nil
}
