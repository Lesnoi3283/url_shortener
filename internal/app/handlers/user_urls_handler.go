package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"go.uber.org/zap"
)

type UserURLsHandler struct {
	URLStorage URLStorageInterface
	Conf       config.Config
	Logger     zap.SugaredLogger
}

type URLData struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

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
		//res.WriteHeader(http.StatusInternalServerError)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	URLDatas := make([]URLData, 0)
	URLsFromDB, err := h.URLStorage.GetUserUrls(req.Context(), userID)
	if err != nil {
		h.Logger.Error("UserURLsHandler error while trying to get user`s urls", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	} else if len(URLsFromDB) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	for _, el := range URLsFromDB {
		URLDatas = append(URLDatas, URLData{
			ShortURL:    h.Conf.BaseAddress + "/" + el.Short,
			OriginalURL: el.Long,
		})
	}

	JSONResp, err := json.Marshal(URLDatas)
	if err != nil {
		h.Logger.Error("UserURLsHandler error while marshalling JSON", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(JSONResp)

}
