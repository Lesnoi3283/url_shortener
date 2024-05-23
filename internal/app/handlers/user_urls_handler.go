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
	DeleteBatchWithUserID(userID int) (urlsChan chan string, err error)
}

type URLData struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (h *UserURLsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	userIDFromContext := req.Context().Value(middlewares.UserIDContextKey)
	userID, ok := (userIDFromContext).(int)
	if (userIDFromContext != nil) || (!ok) {
		h.Logger.Error("UserURLsHandler just got user id `-1` somehow")
		res.WriteHeader(http.StatusInternalServerError)
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
