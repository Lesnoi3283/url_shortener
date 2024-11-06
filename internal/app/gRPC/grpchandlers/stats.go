package grpchandlers

import (
	"context"
	"github.com/Lesnoi3283/url_shortener/internal/app/gRPC/proto"
	"github.com/Lesnoi3283/url_shortener/internal/app/logic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ShortenerServer) Stats(ctx context.Context, req *emptypb.Empty) (*proto.StatsResponse, error) {
	URLs, users, err := logic.GetStats(ctx, s.Storage)
	if err != nil {
		s.Logger.Errorf("Cant get stats, err: %v", err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}
	return &proto.StatsResponse{
		UrlsAmount:  uint64(URLs),
		UsersAmount: uint32(users),
	}, nil
}
