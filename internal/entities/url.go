package entities

type URL struct {
	Real    string
	Short   string
	Storage URLStorageInterface
}

type URLStorageInterface interface {
	Save(URL) error
	Get(string) (URL, error)
	//remove(Real) error
}
