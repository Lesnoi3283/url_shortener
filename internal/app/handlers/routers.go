package handlers

import "github.com/go-chi/chi"

func MyRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/", URLShortenerHandler)
	r.Get("/{url}", ShortURLRedirectHandler)

	return r
}
