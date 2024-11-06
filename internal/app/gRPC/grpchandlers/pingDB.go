package grpchandlers

import (
	"context"
	"github.com/Lesnoi3283/url_shortener/internal/app/logic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ShortenerServer) PingDB(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	err := logic.PingDB(s.Storage)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "Database doesn`t work now")
	} else {
		return &emptypb.Empty{}, nil
	}
}
