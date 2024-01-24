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

	//handlers setting
	r.Post("/", middlewares.LoggerMW(&URLShortener, logger))
	r.Get("/{url}", middlewares.LoggerMW(&shortURLRedirect, logger))

	return r
}
