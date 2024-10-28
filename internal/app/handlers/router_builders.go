package handlers

import (
	"fmt"
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net"
)

// NewRouter builds new chi.Router with handlers. User just have to run it with http.ListenAndServe or something else.
func NewRouter(conf config.Config, store URLStorageInterface, logger zap.SugaredLogger) (chi.Router, error) {
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
	userURLs := UserURLsHandler{
		URLStorage: store,
		Conf:       conf,
		Logger:     logger,
	}
	deleteURLs := DeleteURLsHandler{
		URLStorage: store,
		Conf:       conf,
		Log:        logger,
	}
	pingDB := PingDBHandler{
		DB: store,
	}
	stats := StatsHandler{
		log:     logger,
		storage: store,
	}

	r.Use(middlewares.LoggerMW(logger))
	r.Use(middlewares.CompressionMW(logger))
	r.Use(middlewares.AuthMW(store, logger))

	trustedSubnet := &net.IPNet{}
	if conf.TrustedSubnet != "" {
		var err error
		_, trustedSubnet, err = net.ParseCIDR(conf.TrustedSubnet)
		if err != nil {
			return nil, fmt.Errorf("error parsing trusted subnet: %w", err)
		}
	}
	r.Use(middlewares.SubnetFilterMW(trustedSubnet, logger))

	r.Get("/api/user/urls", userURLs.ServeHTTP)
	r.Post("/", URLShortener.ServeHTTP)
	r.Get("/{url}", shortURLRedirect.ServeHTTP)
	r.Post("/api/shorten", shortener.ServeHTTP)
	r.Post("/api/shorten/batch", shortenBatch.ServeHTTP)
	r.Delete("/api/user/urls", deleteURLs.ServeHTTP)
	r.Get("/ping", pingDB.ServeHTTP)
	r.Get("/api/internal/stats", stats.ServeHTTP)

	return r, nil
}
