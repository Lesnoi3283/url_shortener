package handlers

import (
	"crypto/sha256"
	"fmt"
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/entities"
	"github.com/Lesnoi3283/url_shortener/internal/storages"
	"github.com/Lesnoi3283/url_shortener/pkg/databases/justamap"
	"github.com/go-chi/chi"
	"io"
	"net/http"
)

var jm justamap.JustAMap = justamap.JustAMap{Store: make(map[string]string)}

func ShortURLRedirectHandler(res http.ResponseWriter, req *http.Request) {

	//reading data from request
	shorted := chi.URLParam(req, "url")

	database := jm
	urlStorage := storages.URLStorage{
		DB: &database,
	}
	url := entities.URL{
		Short:   shorted,
		Storage: &urlStorage,
	}

	//reading from db
	url, err := url.Storage.Get(url.Short)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Printf("URL was not found: %v\n", err)
		return
	}

	//response preparing
	res.Header().Set("Location", url.Real)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func URLShortenerHandler(res http.ResponseWriter, req *http.Request) {

	//read request params
	str, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error while reading reqBody")
		return
	}

	//db set
	database := jm
	//todo: Q почему он ругался на `database`, но пропустил `&database`
	urlStorage := storages.URLStorage{
		DB: &database,
	}
	url := entities.URL{
		Real:    string(str),
		Storage: &urlStorage,
	}

	//url shorting
	h := sha256.New()
	h.Write(str)
	strShort := fmt.Sprintf("%x", h.Sum(nil))
	strShort = strShort[:16]
	url.Short = strShort

	//url saving
	err = url.Storage.Save(url)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error while saving to db")
		fmt.Println(err)
		return
	}

	//response making
	res.Header().Set("Content-Type", "text/plain")
	toRet := config.BaseAddress + url.Short
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(toRet))
}
