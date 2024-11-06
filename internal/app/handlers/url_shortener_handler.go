package handlers

import (
	"errors"
	"github.com/Lesnoi3283/url_shortener/internal/app/logic"
	"go.uber.org/zap"
	"io"
	"net/http"

	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"github.com/Lesnoi3283/url_shortener/pkg/databases"
	"github.com/go-chi/chi"
)

// ShortURLRedirectHandler is a handler struct. Use it`s ServeHTTP func.
type ShortURLRedirectHandler struct {
	URLStorage logic.URLStorageInterface
	Log        zap.SugaredLogger
}

// ServeHTTP reads short URL from given URLParam and redirects user to an original URL.
func (h *ShortURLRedirectHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//reading data from request
	shorted := chi.URLParam(req, "url")

	//reading from DB
	fullURL, err := logic.GetOriginalURL(req.Context(), shorted, h.URLStorage)
	if errors.Is(err, databases.ErrURLWasDeleted()) {
		res.WriteHeader(http.StatusGone)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		h.Log.Warnf("error while getting an original URL: %v", err)
		return
	}

	//response preparing
	res.Header().Set("Location", fullURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

// URLShortenerHandler is a handler struct. Use it`s ServeHTTP func.
type URLShortenerHandler struct {
	Conf       config.Config
	URLStorage logic.URLStorageInterface
	Log        zap.SugaredLogger
}

// ServeHTTP shorts a given URL (plain text), saves it in a storage and returns a short version.
func (h *URLShortenerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//this var is necessary. Because it helps to change status code to 409 if url already exists
	successStatus := http.StatusCreated

	//read request params
	realURLBytes, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		h.Log.Errorf("Error while reading reqBody: %v", err)
		return
	}

	//url saving
	userIDFromContext := req.Context().Value(middlewares.UserIDContextKey)
	userID, ok := (userIDFromContext).(int)
	var shortURL string
	if (userIDFromContext != nil) && (ok) {
		shortURL, err = logic.Shorten(req.Context(), realURLBytes, h.Conf.BaseAddress, h.URLStorage, userID)
	} else {
		shortURL, err = logic.Shorten(req.Context(), realURLBytes, h.Conf.BaseAddress, h.URLStorage, -1)
	}

	if err != nil {
		var alrExErr *databases.AlreadyExistsError
		if errors.As(err, &alrExErr) {
			shortURL = alrExErr.ShortURL
			successStatus = http.StatusConflict
		} else {
			res.WriteHeader(http.StatusInternalServerError)
			h.Log.Errorf("Error while shortening URL: %v\n", err)
			return
		}
	}

	//response making
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(successStatus)
	res.Write([]byte(h.Conf.BaseAddress + "/" + shortURL))
}
