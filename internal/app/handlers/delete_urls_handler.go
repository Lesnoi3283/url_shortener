package handlers

import (
	"encoding/json"
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"go.uber.org/zap"
	"net/http"
)

type DeleteURLsHandler struct {
	URLStorage URLStorageInterface
	Conf       config.Config
	Log        zap.SugaredLogger
}

func (h *DeleteURLsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//read request params

	shortURLs := make([]string, 0)

	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&shortURLs)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		h.Log.Error("Error while decoding req body", zap.Error(err))
		return
	}

	userIDFromContext := req.Context().Value(middlewares.UserIDContextKey)
	userID, ok := (userIDFromContext).(int)
	if userIDFromContext == nil || !ok {
		res.WriteHeader(http.StatusInternalServerError)
		h.Log.Error("UserID is nil")
		return
	}

	inputCh, err := h.URLStorage.DeleteBatchWithUserID(userID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		h.Log.Error("Error while deleting URLs", zap.Error(err))
		return
	}
	//fan-out
	go func() {
		defer close(inputCh)
		for _, URL := range shortURLs {
			inputCh <- URL
		}
	}()

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusAccepted)
}
