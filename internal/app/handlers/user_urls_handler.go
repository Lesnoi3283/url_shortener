package handlers

import (
	"encoding/json"
	"github.com/Lesnoi3283/url_shortener/internal/app/logic"
	"net/http"

	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"go.uber.org/zap"
)

// UserURLsHandler is a handler struct. Use it`s ServeHTTP func.
type UserURLsHandler struct {
	URLStorage logic.URLStorageInterface
	Conf       config.Config
	Logger     zap.SugaredLogger
}

// ServeHTTP returns a JSON array with all users`s urls.
func (h *UserURLsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	cookie, err := req.Cookie(middlewares.JwtCookieName)
	if err != nil {
		h.Logger.Error("UserURLsHandler cookie get err", zap.Error(err))
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	//CTX is not needed! Because user MUST be authorised BEFORE this request! GitHub tests sayed me that...
	token := cookie.Value
	userID := middlewares.GetUserID(token)
	if userID == -1 {
		h.Logger.Error("UserURLsHandler just got user id `-1` somehow. Probably JWT is not valid")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	//get URLs
	usersURLs, err := logic.GetUsersURLs(req.Context(), h.URLStorage, h.Conf.BaseAddress, userID)
	if err != nil {
		h.Logger.Error("UserURLsHandler get err", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(usersURLs) == 0 {
		res.WriteHeader(http.StatusNoContent)
		h.Logger.Debugf("users urls is empty for user with ID = `%v`", userID)
		return
	}

	//make a response
	JSONResp, err := json.Marshal(usersURLs)
	if err != nil {
		h.Logger.Error("UserURLsHandler error while marshalling JSON", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(JSONResp)
}
