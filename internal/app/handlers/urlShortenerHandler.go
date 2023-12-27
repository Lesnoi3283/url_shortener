package handlers

import (
	"github.com/Lesnoi3283/url_shortener/internal/entities"
	"github.com/Lesnoi3283/url_shortener/internal/storages"
	"net/http"
)

func URLShortenerHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		str := "ya.ru" //imagine we`ve got it from req

		var urlStorage storages.UrlStorage
		var url entities.Url = entities.Url{
			Real:    str,
			Storage: urlStorage,
		}

		strShort := "short" //imagine we`ve hashed it or smg i also dunno
		url.Short = strShort
		err := url.Storage.Save(url)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "text/plain")
		res.Write([]byte(url.Short))
		return

	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}
