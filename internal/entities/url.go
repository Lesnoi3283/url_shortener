package entities

type Url struct {
	Real    string
	Short   string
	Storage UrlStorage
}

type UrlStorage interface {
	Save(Url) error
	Get(Url) (Url, error)
	//remove(Real) error
}
