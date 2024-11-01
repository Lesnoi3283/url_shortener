package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Lesnoi3283/url_shortener/internal/app/logic"
	"github.com/Lesnoi3283/url_shortener/pkg/databases"
	"io"
	"log"
	"net/http"

	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"go.uber.org/zap"
)

// ShortenHandler is a handler struct. Use it`s ServeHTTP func.
type ShortenHandler struct {
	URLStorage logic.URLStorageInterface
	Conf       config.Config
	Log        zap.SugaredLogger
}

// ServeHTTP shorts given url (JSON), saves it in a storage and return a short version.
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

	//get userID and short the URL
	userIDFromContext := req.Context().Value(middlewares.UserIDContextKey)
	userID, ok := (userIDFromContext).(int)
	var urlShort string
	if ok {
		urlShort, err = logic.Shorten(req.Context(), []byte(realURL.Val), h.Conf.BaseAddress, h.URLStorage, userID)
	} else {
		urlShort, err = logic.Shorten(req.Context(), []byte(realURL.Val), h.Conf.BaseAddress, h.URLStorage, -1)
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		h.Log.Errorf("Error while shortening URL '%s': %v", realURL.Val, err)
		return
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
		Result: urlShort,
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
