package handlers

import (
	"encoding/json"
	"github.com/Lesnoi3283/url_shortener/internal/app/logic"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"io"
	"log"
	"net/http"

	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
	"go.uber.org/zap"
)

// ShortenBatchHandler is a handler struct. Use it`s ServeHTTP func.
type ShortenBatchHandler struct {
	URLStorage logic.URLStorageInterface
	Conf       config.Config
	Log        zap.SugaredLogger
}

// ServeHTTP shorts all given URLS (in JSON) and saves them in a storage.
// Returns a JSON array with short versions of given URLs.
func (h *ShortenBatchHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//read request params
	bodyBytes, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		h.Log.Error("Error while reading reqBody")
		return
	}

	URLs := make([]entities.URL, 0)

	err = json.Unmarshal(bodyBytes, &URLs)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		h.Log.Error("Error during unmarshalling JSON")
		return
	}

	//get userID and short URLs
	userIDFromContext := req.Context().Value(middlewares.UserIDContextKey)
	userID, ok := (userIDFromContext).(int)
	if ok {
		URLs, err = logic.ShortenBatch(req.Context(), URLs, h.Conf.BaseAddress, h.URLStorage, userID)
	} else {
		URLs, err = logic.ShortenBatch(req.Context(), URLs, h.Conf.BaseAddress, h.URLStorage, -1)
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		h.Log.Errorf("Error while shortening batch of URLs: %v", err)
		return
	}

	//response making
	jsonResponse, err := json.Marshal(URLs)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Default().Println("Error during marshalling JSON response")
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(jsonResponse)
}
