package handlers

import (
	"context"
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// Я точно правильно передаю конф, логгер, бдшки? Слишком много аргументов у функции
func BuildRouter(conf config.Config, store URLStorageInterface, logger zap.SugaredLogger, db DBInterface) chi.Router {
	r := chi.NewRouter()

	//где лучше создавать контекст?
	//handlers building
	URLShortener := URLShortenerHandler{
		ctx:        context.Background(),
		Conf:       conf,
		URLStorage: store,
	}
	shortURLRedirect := ShortURLRedirectHandler{
		ctx:        context.Background(),
		URLStorage: store,
	}
	shortener := shortenHandler{
		ctx:        context.Background(),
		Conf:       conf,
		URLStorage: store,
	}
	pingDB := pingDBHandler{
		db: db,
	}
	shortenBatch := shortenBatchHandler{
		ctx:        context.Background(),
		URLStorage: store,
		Conf:       conf,
	}

	//handlers setting
	r.Post("/", middlewares.LoggerMW(middlewares.CompressionMW(&URLShortener, logger), logger))
	r.Get("/{url}", middlewares.LoggerMW(middlewares.CompressionMW(&shortURLRedirect, logger), logger))
	r.Post("/api/shorten", middlewares.LoggerMW(middlewares.CompressionMW(&shortener, logger), logger))
	r.Get("/ping", middlewares.LoggerMW(middlewares.CompressionMW(&pingDB, logger), logger))
	r.Post("/api/shorten/batch", middlewares.LoggerMW(middlewares.CompressionMW(&shortenBatch, logger), logger))

	return r
}
