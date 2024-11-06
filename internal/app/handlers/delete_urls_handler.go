// Package handlers includes all handlers for the url_shortener server.
package handlers

import (
	"encoding/json"
	"github.com/Lesnoi3283/url_shortener/internal/app/logic"
	"net/http"

	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"go.uber.org/zap"
)

// DeleteURLsHandler is a handler struct. Use it`s ServeHTTP func.
type DeleteURLsHandler struct {
	URLStorage logic.URLStorageInterface
	Conf       config.Config
	Log        zap.SugaredLogger
}

// ServeHTTP deletes all given URLs (in JSON). Only for authorised users.
// If given URL was created by different user - nothing would be deleted.
func (h *DeleteURLsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//read request params

	shortURLs := make([]string, 0)

	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&shortURLs)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		h.Log.Error("Error while decoding req body", zap.Error(err))
		return
	}

	userIDFromContext := req.Context().Value(middlewares.UserIDContextKey)
	userID, ok := (userIDFromContext).(int)
	if userIDFromContext == nil || !ok {
		res.WriteHeader(http.StatusUnauthorized)
		h.Log.Error("UserID is nil")
		return
	}

	err = logic.DeleteURLs(userID, shortURLs, h.URLStorage)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		h.Log.Error("Error while deleting urls", zap.Error(err))
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusAccepted)
}
