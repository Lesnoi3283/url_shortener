package storages

import (
	"github.com/Lesnoi3283/url_shortener/internal/entities"
)

type UrlStorage struct {
	Db DB
}

type DB interface {
	Save(key string, val string) error
	Get(key string) (string, error)
	//mb remove but i dunno if it is necessary
}

func (u *UrlStorage) Save(url entities.Url) error {
	err := u.Db.Save(url.Short, url.Real)
	return err
}

func (u *UrlStorage) Get(short string) (url entities.Url, err error) {
	str, err := u.Db.Get(short)
	url.Real = str
	url.Short = short
	return url, err
}
