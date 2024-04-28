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
		Log:        logger,
	}
	shortenBatch := ShortenBatchHandler{
		URLStorage: store,
		Conf:       conf,
		Log:        logger,
	}

	//handlers setting
	userURLsSotrange, ok := (store).(UserUrlsStorageInterface)
	if ok {
		//normal db service
		userURLs := UserURLsHandler{
			URLStorage: userURLsSotrange,
			Conf:       config.Config{},
			Logger:     logger,
		}
		r.Get("/api/user/urls", middlewares.LoggerMW(middlewares.CompressionMW(middlewares.AuthMW(&userURLs, userURLsSotrange, logger), logger), logger))
		r.Post("/", middlewares.LoggerMW(middlewares.CompressionMW(middlewares.AuthMW(http.HandlerFunc(URLShortener.ServeHTTP), userURLsSotrange, logger), logger), logger)) //вот так надо
		r.Get("/{url}", middlewares.LoggerMW(middlewares.CompressionMW(middlewares.AuthMW(&shortURLRedirect, userURLsSotrange, logger), logger), logger))
		r.Post("/api/shorten", middlewares.LoggerMW(middlewares.CompressionMW(middlewares.AuthMW(&shortener, userURLsSotrange, logger), logger), logger))
		r.Post("/api/shorten/batch", middlewares.LoggerMW(middlewares.CompressionMW(middlewares.AuthMW(&shortenBatch, userURLsSotrange, logger), logger), logger))
	} else {
		//shitty db edition
		r.Post("/", middlewares.LoggerMW(middlewares.CompressionMW(http.HandlerFunc(URLShortener.ServeHTTP), logger), logger)) //вот так надо
		r.Get("/{url}", middlewares.LoggerMW(middlewares.CompressionMW(&shortURLRedirect, logger), logger))
		r.Post("/api/shorten", middlewares.LoggerMW(middlewares.CompressionMW(&shortener, logger), logger))
		r.Post("/api/shorten/batch", middlewares.LoggerMW(middlewares.CompressionMW(&shortenBatch, logger), logger))
	}

	postgresqlDB, ok := (store).(*databases.Postgresql)
	if ok {
		pingDB := PingDBHandler{
			DB: postgresqlDB,
		}
		r.Get("/ping", middlewares.LoggerMW(middlewares.CompressionMW(&pingDB, logger), logger))
	}

	return r
}
