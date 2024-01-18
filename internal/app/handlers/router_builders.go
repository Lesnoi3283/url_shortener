package handlers

import (
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/go-chi/chi"
)

func BuildRouter(conf config.Config, store URLStorageInterface) chi.Router {
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
	r.Post("/", URLShortener.ServeHTTP)
	r.Get("/{url}", shortURLRedirect.ServeHTTP)

	return r
}
