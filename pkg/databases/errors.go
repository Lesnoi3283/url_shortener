package databases

type AlreadyExistsError struct {
	ShortURL string
}

func (a *AlreadyExistsError) Error() string {
	return a.ShortURL
}

func NewAlreadyExistsError(shortURL string) *AlreadyExistsError {
	return &AlreadyExistsError{ShortURL: shortURL}
}

func (a *AlreadyExistsError) Is(target error) bool {
	_, ok := target.(*AlreadyExistsError)
	return ok
}
