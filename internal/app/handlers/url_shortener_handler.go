package handlers

import (
	"crypto/sha256"
	"fmt"
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/entities"
	"github.com/go-chi/chi"
	"io"
	"log"
	"net/http"
)

type ShortURLRedirectHandler struct {
	URLStorage entities.URLStorageInterface
}

func (h *ShortURLRedirectHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//reading data from request
	shorted := chi.URLParam(req, "url")

	//reading from db
	url, err := h.URLStorage.Get(shorted)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		log.Default().Printf("url was not found: %v\n", err)
		return
	}

	//response preparing
	res.Header().Set("Location", url.Real)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

type URLShortenerHandler struct {
	Conf       config.Config
	URLStorage entities.URLStorageInterface
}

func (h *URLShortenerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//read request params
	str, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Default().Println("Error while reading reqBody")
		return
	}

	url := entities.URL{
		Real: string(str),
	}

	//url shorting
	hasher := sha256.New()
	hasher.Write(str)
	strShort := fmt.Sprintf("%x", hasher.Sum(nil))
	strShort = strShort[:16]
	url.Short = strShort

	//url saving
	err = h.URLStorage.Save(url)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Default().Println("Error while saving to db")
		log.Default().Println(err)
		return
	}

	//response making
	res.Header().Set("Content-Type", "text/plain")
	toRet := h.Conf.BaseAddress + "/" + url.Short
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(toRet))
}
