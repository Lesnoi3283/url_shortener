package handlers

import (
	"context"
	"encoding/json"
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"go.uber.org/zap"
	"net/http"
)

type UserURLsHandler struct {
	URLStorage UserUrlsStorageInterface
	Conf       config.Config
	Logger     zap.SugaredLogger
}

type UserUrlsStorageInterface interface {
	URLStorageInterface
	GetUserUrls(ctx context.Context, userID int) ([]struct {
		Long  string
		Short string
	}, error)
	SaveWithUserID(ctx context.Context, userID int, short string, full string) error
	SaveBatchWithUserID(ctx context.Context, userID int, urls []entities.URL) error
	CreateUser(ctx context.Context) (int, error)
}

type URLData struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (h *UserURLsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	userIDFromContext := req.Context().Value(middlewares.UserIDContextName)
	userID, ok := (userIDFromContext).(int)

	h.Logger.Info("shit", zap.Int("userid", userID))
	if userIDFromContext == nil {
		h.Logger.Error("UserURLsHandler cookie get err")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !ok {
		h.Logger.Error("UserURLsHandler something wasn`t ok when i tried to cast userIDFromContext to int")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userID == -1 {
		h.Logger.Error("UserURLsHandler just got user id `-1` somehow")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	URLDatas := make([]URLData, 5)
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
			ShortURL:    el.Short,
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
	res.Write(JSONResp)
	res.WriteHeader(http.StatusOK)
}
