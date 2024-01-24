package handlers

import (
	"crypto/sha256"
	"fmt"
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/go-chi/chi"
	"io"
	"log"
	"net/http"
)

type URLStorageInterface interface {
	Save(short string, full string) error
	Get(short string) (full string, err error)
	//remove(Real) error
}

type ShortURLRedirectHandler struct {
	URLStorage URLStorageInterface
}

func (h *ShortURLRedirectHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//reading data from request
	shorted := chi.URLParam(req, "url")

	//reading from db
	fullURL, err := h.URLStorage.Get(shorted)
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
	err = h.URLStorage.Save(urlShort, realURL)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Default().Println("Error while saving to db")
		log.Default().Println(err)
		return
	}

	//response making
	res.Header().Set("Content-Type", "text/plain")
	toRet := h.Conf.BaseAddress + "/" + urlShort
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(toRet))
}
