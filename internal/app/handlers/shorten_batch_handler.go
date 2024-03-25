package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/Lesnoi3283/url_shortener/config"
	"io"
	"log"
	"net/http"
	"strconv"
)

type shortenBatchHandler struct {
	ctx        context.Context
	URLStorage URLStorageInterface
	Conf       config.Config
}

type URL struct {
	Short string
	Long  string
}

func (h *shortenBatchHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//read request params
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Default().Println("Error while reading reqBody")
		return
	}

	type URLGot struct {
		CorrelationId string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	URLsGot := make([]URLGot, 0)

	err = json.Unmarshal(bodyBytes, &URLsGot)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Default().Println("Error during unmarshalling JSON")
		return
	}

	type URLShorten struct {
		CorrelationId string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	URLsToSave := make([]URL, 0)
	URLsToReturn := make([]URLShorten, 0)

	for i, url := range URLsGot {
		//url shorting
		hasher := sha256.New()
		hasher.Write([]byte(url.OriginalURL))
		urlShort := fmt.Sprintf("%x", hasher.Sum(nil))
		urlShort = urlShort[:16]

		URLsToSave = append(URLsToSave, URL{
			Short: urlShort,
			Long:  url.OriginalURL,
		})
		URLsToReturn = append(URLsToReturn, URLShorten{
			CorrelationId: strconv.Itoa(i),
			ShortURL:      urlShort,
		})
	}

	//url saving
	err = h.URLStorage.SaveBatch(h.ctx, URLsToSave)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Default().Println("Error while saving to db")
		log.Default().Println(err)
		return
	}

	//response making
	jsonResponce, err := json.Marshal(URLsToReturn)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Default().Println("Error during marshalling JSON responce")
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(jsonResponce)
}
