package handlers

import (
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"github.com/Lesnoi3283/url_shortener/pkg/databases"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
)

func NewRouter(conf config.Config, store URLStorageInterface, logger zap.SugaredLogger, db DBInterface) chi.Router {
	r := chi.NewRouter()

	//handlers building
	URLShortener := URLShortenerHandler{
		Conf:       conf,
		URLStorage: store,
	}
	shortURLRedirect := ShortURLRedirectHandler{
		URLStorage: store,
	}
	shortener := ShortenHandler{
		Conf:       conf,
		URLStorage: store,
	}
	shortenBatch := ShortenBatchHandler{
		URLStorage: store,
		Conf:       conf,
	}

	//handlers setting
	r.Post("/", middlewares.LoggerMW(middlewares.CompressionMW(http.HandlerFunc(URLShortener.ServeHTTP), logger), logger)) //вот так надо
	r.Get("/{url}", middlewares.LoggerMW(middlewares.CompressionMW(&shortURLRedirect, logger), logger))
	r.Post("/api/shorten", middlewares.LoggerMW(middlewares.CompressionMW(&shortener, logger), logger))
	r.Post("/api/shorten/batch", middlewares.LoggerMW(middlewares.CompressionMW(&shortenBatch, logger), logger))
	postgresqlDB, ok := (store).(*databases.Postgresql)
	if ok {
		pingDB := PingDBHandler{
			DB: postgresqlDB,
		}
		r.Get("/ping", middlewares.LoggerMW(middlewares.CompressionMW(&pingDB, logger), logger))
	}
	//r.Get("/ping", middlewares.LoggerMW(middlewares.CompressionMW(&pingDB, logger), logger))

	return r
}
