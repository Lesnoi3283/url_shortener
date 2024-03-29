package handlers

import (
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func BuildRouter(conf config.Config, store URLStorageInterface, logger zap.SugaredLogger) chi.Router {
	r := chi.NewRouter()

	//handlers building
	URLShortener := URLShortenerHandler{
		Conf:       conf,
		URLStorage: store,
	}
	shortURLRedirect := ShortURLRedirectHandler{
		URLStorage: store,
	}
	shortener := shortenHandler{
		Conf:       conf,
		URLStorage: store,
	}

	//handlers setting
	r.Post("/", middlewares.LoggerMW(middlewares.CompressionMW(&URLShortener, logger), logger))
	r.Get("/{url}", middlewares.LoggerMW(middlewares.CompressionMW(&shortURLRedirect, logger), logger))
	r.Post("/api/shorten", middlewares.LoggerMW(middlewares.CompressionMW(&shortener, logger), logger))

	return r
}
