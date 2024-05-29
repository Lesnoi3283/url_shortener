package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
)

type ShortenBatchHandler struct {
	URLStorage URLStorageInterface
	Conf       config.Config
	Log        zap.SugaredLogger
}

func (h *ShortenBatchHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//read request params
	bodyBytes, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		h.Log.Error("Error while reading reqBody")
		return
	}

	type URLGot struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	URLsGot := make([]URLGot, 0)

	err = json.Unmarshal(bodyBytes, &URLsGot)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		h.Log.Error("Error during unmarshalling JSON")
		return
	}

	type URLShorten struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	URLsToSave := make([]entities.URL, 0)
	URLsToReturn := make([]URLShorten, 0)

	for i, url := range URLsGot {
		//url shorting
		hasher := sha256.New()
		hasher.Write([]byte(url.OriginalURL))
		urlShort := fmt.Sprintf("%x", hasher.Sum(nil))
		urlShort = urlShort[:16]

		URLsToSave = append(URLsToSave, entities.URL{
			Short: urlShort,
			Long:  url.OriginalURL,
		})
		URLsToReturn = append(URLsToReturn, URLShorten{
			CorrelationID: URLsGot[i].CorrelationID,
			ShortURL:      h.Conf.BaseAddress + "/" + urlShort,
		})
	}

	//url saving
	userIDFromContext := req.Context().Value(middlewares.UserIDContextKey)
	userID, ok := (userIDFromContext).(int)
	if (userIDFromContext != nil) && (ok) {
		err = h.URLStorage.SaveBatchWithUserID(req.Context(), userID, URLsToSave)
	} else {
		err = h.URLStorage.SaveBatch(req.Context(), URLsToSave)
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		h.Log.Error("Error while saving to DB", zap.Error(err))
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
