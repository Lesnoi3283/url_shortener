package grpchandlers

import (
	"context"
	"github.com/Lesnoi3283/url_shortener/internal/app/gRPC/proto"
	"github.com/Lesnoi3283/url_shortener/internal/app/logic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ShortenerServer) GetOriginalURL(ctx context.Context, req *proto.GetOriginalURLRequest) (*proto.GetAnOriginalURLResponse, error) {
	url, err := logic.GetOriginalURL(ctx, req.ShortUrl, s.Storage)
	if err != nil {
		s.Logger.Debugf("Original URL not found. Given short: %v", url)
		return nil, status.Error(codes.NotFound, err.Error())
	}
	res := &proto.GetAnOriginalURLResponse{
		Url: url,
	}
	return res, nil
}
