package entities

type URL struct {
	Real  string
	Short string
}

type URLStorageInterface interface {
	Save(URL) error
	Get(string) (URL, error)
	//remove(Real) error
}
