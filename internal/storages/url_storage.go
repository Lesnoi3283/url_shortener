package storages

import (
	"github.com/Lesnoi3283/url_shortener/internal/entities"
)

type URL struct {
	Real  string
	Short string
}

type URLStorageInterface interface {
	Save(URL) error
	Get(string) (URL, error)
	//remove(Real) error
}

type URLStorage struct {
	DB DBInterface
}

type DBInterface interface {
	Save(key string, val string) error
	Get(key string) (string, error)
	//mb remove but i dunno if it is necessary
}

func (u *URLStorage) Save(url entities.URL) error {
	err := u.DB.Save(url.Short, url.Real)
	return err
}

func (u *URLStorage) Get(short string) (url entities.URL, err error) {
	str, err := u.DB.Get(short)
	url.Real = str
	url.Short = short
	return url, err
}
