package handlers

import (
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"github.com/Lesnoi3283/url_shortener/pkg/databases"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"time"
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
	//todo: добавить лимитер мв (ради эксперимента)
	userURLsSotrange, ok := (store).(UserUrlsStorageInterface)
	if ok {
		//normal db service
		userURLs := UserURLsHandler{
			URLStorage: userURLsSotrange,
			Conf:       conf,
			Logger:     logger,
		}
		deleteURLs := DeleteURLsHandler{
			URLStorage: userURLsSotrange,
			Conf:       conf,
			Log:        logger,
		}

		r.Use(middlewares.LoggerMW(logger))
		requestManager := middlewares.NewRequestManager(100, time.Minute)
		r.Use(middlewares.RequestLimiterMW(logger, requestManager))
		r.Use(middlewares.CompressionMW(logger))
		r.Use(middlewares.AuthMW(userURLsSotrange, logger))

		r.Get("/api/user/urls", userURLs.ServeHTTP)
		r.Post("/", URLShortener.ServeHTTP)
		r.Get("/{url}", shortURLRedirect.ServeHTTP)
		r.Post("/api/shorten", shortener.ServeHTTP)
		r.Post("/api/shorten/batch", shortenBatch.ServeHTTP)
		r.Delete("/api/user/urls", deleteURLs.ServeHTTP)
	} else {
		//shitty db edition
		r.Post("/", URLShortener.ServeHTTP)
		r.Get("/{url}", shortURLRedirect.ServeHTTP)
		r.Post("/api/shorten", shortener.ServeHTTP)
		r.Post("/api/shorten/batch", shortenBatch.ServeHTTP)
	}

	postgresqlDB, ok := (store).(*databases.Postgresql)
	if ok {
		pingDB := PingDBHandler{
			DB: postgresqlDB,
		}
		r.Get("/ping", pingDB.ServeHTTP)
	}

	return r
}
