package handlers

import (
	"context"
	"github.com/Lesnoi3283/url_shortener/config"
	"net/http"
)

type UserURLsHandler struct {
	URLStorage URLStorageInterface
	Conf       config.Config
}

type UserUrlsStorageInterface interface {
	URLStorageInterface
	GetUserUrls(ctx context.Context, userID int) ([]struct {
		Long  string
		Short string
	}, error)
	CreateUser(ctx context.Context) (int, error)
}

func (h *UserURLsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

}
