package storages

import (
	"github.com/Lesnoi3283/url_shortener/internal/entities"
)

type UrlStorage struct {
	db DB
}

type DB interface {
	Save(key string, val string) error
	Get(key string) (string, error)
	//mb remove but i dunno if it is necessary
}

func (u UrlStorage) Save(url entities.Url) error {
	err := u.db.Save(url.Short, url.Real)
	return err
}

func (u UrlStorage) Get(url entities.Url) (entities.Url, error) {
	str, err := u.db.Get(url.Short)
	url.Real = str
	return url, err
}
