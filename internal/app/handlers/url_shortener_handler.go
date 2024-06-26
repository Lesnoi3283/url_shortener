package handlers

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"github.com/Lesnoi3283/url_shortener/pkg/databases"
	"github.com/go-chi/chi"
	"io"
	"log"
	"net/http"
)

//go:generate mockgen -source=url_shortener_handler.go -destination=mocks/mock_DBInterface.go -package=mocks github.com/Lesnoi3283/url_shortener/internal/app/handlers URLStorageInterface

type URLStorageInterface interface {
	Save(ctx context.Context, url entities.URL) error
	SaveBatch(ctx context.Context, urls []entities.URL) error
	Get(ctx context.Context, short string) (full string, err error)
	SaveWithUserID(ctx context.Context, userID int, url entities.URL) error
	SaveBatchWithUserID(ctx context.Context, userID int, urls []entities.URL) error
	DeleteBatchWithUserID(userID int) (urlsChan chan string, err error)
	GetUserUrls(ctx context.Context, userID int) ([]entities.URL, error)
	Ping() error
	CreateUser(ctx context.Context) (int, error)
}

type ShortURLRedirectHandler struct {
	URLStorage URLStorageInterface
}

func (h *ShortURLRedirectHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//reading data from request
	shorted := chi.URLParam(req, "url")

	//reading from DB
	fullURL, err := h.URLStorage.Get(req.Context(), shorted)
	if errors.Is(err, databases.ErrURLWasDeleted()) {
		res.WriteHeader(http.StatusGone)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("fullURL was not found: %v\n", err)
		return
	}

	//response preparing
	res.Header().Set("Location", fullURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

type URLShortenerHandler struct {
	Conf       config.Config
	URLStorage URLStorageInterface
}

func (h *URLShortenerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//this var is necessary. Because it helps to change status code to 409 if url already exists
	successStatus := http.StatusCreated

	//read request params
	str, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Default().Println("Error while reading reqBody")
		return
	}

	realURL := string(str)

	//url shorting
	hasher := sha256.New()
	hasher.Write(str)
	urlShort := fmt.Sprintf("%x", hasher.Sum(nil))
	urlShort = urlShort[:16]

	//url saving
	userIDFromContext := req.Context().Value(middlewares.UserIDContextKey)
	userID, ok := (userIDFromContext).(int)
	if (userIDFromContext != nil) && (ok) {
		err = h.URLStorage.SaveWithUserID(req.Context(), userID, entities.URL{
			Short: urlShort,
			Long:  realURL,
		})
	} else {
		err = h.URLStorage.Save(req.Context(), entities.URL{
			Short: urlShort,
			Long:  realURL,
		})
	}

	var alrExErr *databases.AlreadyExistsError
	if errors.As(err, &alrExErr) {
		urlShort = alrExErr.ShortURL
		successStatus = http.StatusConflict
	} else if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Default().Println("Error while saving to DB")
		log.Default().Println(err)
		return
	}

	//response making
	res.Header().Set("Content-Type", "text/plain")
	toRet := h.Conf.BaseAddress + "/" + urlShort
	res.WriteHeader(successStatus)
	res.Write([]byte(toRet))
}
