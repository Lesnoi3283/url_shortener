package logic

import (
	"context"
	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
)

//go:generate mockgen -source=url_shortener_handler.go -destination=mocks/mock_DBInterface.go -package=mocks github.com/Lesnoi3283/url_shortener/internal/app/handlers URLStorageInterface

// URLStorageInterface is a main database interface.
type URLStorageInterface interface {
	Save(ctx context.Context, url entities.URL) error
	SaveBatch(ctx context.Context, urls []entities.URL) error
	Get(ctx context.Context, short string) (full string, err error)
	SaveWithUserID(ctx context.Context, userID int, url entities.URL) error
	SaveBatchWithUserID(ctx context.Context, userID int, urls []entities.URL) error
	DeleteBatchWithUserID(userID int) (urlsChan chan string, err error)
	GetUserUrls(ctx context.Context, userID int) ([]entities.URL, error)
	Ping() error
	CreateUser(ctx context.Context) (int, error)
	GetUsersCount(ctx context.Context) (int, error)
	GetShortURLCount(ctx context.Context) (int, error)
}
