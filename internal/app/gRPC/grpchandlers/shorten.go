package grpchandlers

import (
	"context"
	"errors"
	"github.com/Lesnoi3283/url_shortener/internal/app/gRPC/interceptors"
	"github.com/Lesnoi3283/url_shortener/internal/app/gRPC/proto"
	"github.com/Lesnoi3283/url_shortener/internal/app/logic"
	"github.com/Lesnoi3283/url_shortener/pkg/databases"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ShortenerServer) Shorten(ctx context.Context, req *proto.ShortenRequest) (*proto.ShortenResponse, error) {
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

	//shorten
	short, err := logic.Shorten(ctx, []byte(req.OriginalUrl), s.Conf.BaseAddress, s.Storage, userIDInt)
	alrExistsErr := &databases.AlreadyExistsError{}
	if errors.As(err, &alrExistsErr) {
		short = alrExistsErr.ShortURL
		return &proto.ShortenResponse{Shorten: short}, status.Error(codes.AlreadyExists, "Already exists")
	}
	if err != nil {
		s.Logger.Errorf("Shorten err: %v", err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}
	return &proto.ShortenResponse{Shorten: short}, nil
}
