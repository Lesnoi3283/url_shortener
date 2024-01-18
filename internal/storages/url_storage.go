package storages

type URL struct {
	Real  string
	Short string
}

type URLStorage struct {
	DB DBInterface
}

type DBInterface interface {
	Save(key string, val string) error
	Get(key string) (string, error)
	//mb remove but i dunno if it is necessary
}

func (u *URLStorage) Save(url URL) error {
	err := u.DB.Save(url.Short, url.Real)
	return err
}

func (u *URLStorage) Get(short string) (url URL, err error) {
	str, err := u.DB.Get(short)
	url.Real = str
	url.Short = short
	return url, err
}
