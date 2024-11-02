package grpchandlers

import (
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/gRPC/proto"
	"github.com/Lesnoi3283/url_shortener/internal/app/logic"
	"go.uber.org/zap"
)

type ShortenerServer struct {
	proto.UnimplementedURLShortenerServiceServer
	Storage logic.URLStorageInterface
	Logger  zap.SugaredLogger
	Conf    *config.Config
}
