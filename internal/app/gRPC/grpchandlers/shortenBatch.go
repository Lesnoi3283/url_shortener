package grpchandlers

import (
	"context"
	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
	"github.com/Lesnoi3283/url_shortener/internal/app/gRPC/interceptors"
	"github.com/Lesnoi3283/url_shortener/internal/app/gRPC/proto"
	"github.com/Lesnoi3283/url_shortener/internal/app/logic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ShortenerServer) ShortenBatch(ctx context.Context, req *proto.ShortenBatchRequest) (*proto.ShortenBatchResponse, error) {
	//auth
	userID := ctx.Value(interceptors.UserIDContextKey)
	userIDInt := -1
	if userID == nil || userID == -1 {
		s.Logger.Debug("UserID not found req ctx, will use `-1`")
	} else {
		var ok bool
		userIDInt, ok = userID.(int)
		if !ok {
			s.Logger.Warnf("User ID is not an int, real type: `%T`, value: `%v`", userID, userID)
			return nil, status.Error(codes.InvalidArgument, "User ID is not an int")
		}
	}

	//parse request
	URLs := make([]entities.URL, len(req.Urls))
	for i, url := range req.Urls {
		URLs[i] = entities.URL{
			CorrelationID: url.CorrelationId,
			OriginalURL:   url.OriginalUrl,
		}
	}

	//shorten
	URLs, err := logic.ShortenBatch(ctx, URLs, s.Conf.BaseAddress, s.Storage, userIDInt)
	if err != nil {
		s.Logger.Errorf("ShortenBatch error: %v", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	//prepare response
	respURLs := make([]*proto.ShortenBatchResponse_URL, len(URLs))
	for i, url := range URLs {
		respURLs[i] = &proto.ShortenBatchResponse_URL{
			CorrelationId: url.CorrelationID,
			ShortenUrl:    url.ShortURL,
		}
	}

	//return the response
	response := &proto.ShortenBatchResponse{
		Urls: respURLs,
	}
	return response, nil
}
