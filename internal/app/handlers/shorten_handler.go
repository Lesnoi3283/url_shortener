package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"github.com/Lesnoi3283/url_shortener/pkg/databases"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
)

type ShortenHandler struct {
	URLStorage URLStorageInterface
	Conf       config.Config
	Log        zap.SugaredLogger
}

func (h *ShortenHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//this var is using for changing status to 409 if url already exists
	successStatus := http.StatusCreated

	//read request params
	bodyBytes, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		h.Log.Error("Error while reading req body", zap.Error(err))
		return
	}

	//unmarshalling JSON
	realURL := struct {
		Val string `json:"url"`
	}{}

	err = json.Unmarshal(bodyBytes, &realURL)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Default().Println("Error during unmarshalling JSON")
		return
	}

	//url shorting
	hasher := sha256.New()
	hasher.Write(bodyBytes)
	urlShort := fmt.Sprintf("%x", hasher.Sum(nil))
	urlShort = urlShort[:16]

	//url saving
	userURLsStorage, storageOk := (h.URLStorage).(UserUrlsStorageInterface)
	userIDFromContext := req.Context().Value(middlewares.UserIDContextKey)
	userID, ok := (userIDFromContext).(int)
	if (userIDFromContext != nil) && (ok) && storageOk {
		err = userURLsStorage.SaveWithUserID(req.Context(), userID, urlShort, realURL.Val)
	} else {
		err = h.URLStorage.Save(req.Context(), urlShort, realURL.Val)
	}
	var alrExErr *databases.AlreadyExistsError
	if errors.As(err, &alrExErr) {
		urlShort = alrExErr.ShortURL
		successStatus = http.StatusConflict
	} else if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		h.Log.Error("Error while saving to DB", zap.Error(err))
		return
	}

	//response making
	responce := struct {
		Result string `json:"result"`
	}{
		Result: h.Conf.BaseAddress + "/" + urlShort,
	}

	jsonResponce, err := json.Marshal(responce)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Default().Println("Error during marshalling JSON responce")
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(successStatus)
	res.Write(jsonResponce)

}
