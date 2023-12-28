package handlers

import (
	"crypto/sha256"
	"fmt"
	"github.com/Lesnoi3283/url_shortener/internal/entities"
	"github.com/Lesnoi3283/url_shortener/internal/storages"
	"github.com/Lesnoi3283/url_shortener/pkg/databases/justamap"
	"io"
	"net/http"
	"strings"
)

var jm justamap.JustAMap = justamap.JustAMap{Store: make(map[string]string)}

func URLShortenerHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {

		//read request params
		str, err := io.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Println("Error while reading body")
			return
		}

		//db set
		database := jm
		//todo: Q почему он ругался на `database`, но пропустил `&database`
		var urlStorage storages.UrlStorage = storages.UrlStorage{
			Db: &database,
		}
		var url entities.Url = entities.Url{
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
		toRet := "http://localhost:8080/" + url.Short
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(toRet))

		return
	} else if req.Method == http.MethodGet {

		//reading data from request
		shorted := strings.Split(req.URL.String(), "/")
		if len(shorted) != 2 {
			res.WriteHeader(http.StatusBadRequest)
			fmt.Println("Wrong request format")
			return
		}

		//db set
		database := jm
		var urlStorage storages.UrlStorage = storages.UrlStorage{
			Db: &database,
		}
		var url entities.Url = entities.Url{
			Short:   shorted[1],
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

	} else {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}
