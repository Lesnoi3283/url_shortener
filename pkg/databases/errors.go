package databases

type AlreadyExistsError struct {
	ShortUrl string
}

func (a *AlreadyExistsError) Error() string {
	return a.ShortUrl
}

func NewAlreadyExistsError(shortUrl string) *AlreadyExistsError {
	return &AlreadyExistsError{ShortUrl: shortUrl}
}

func (a *AlreadyExistsError) Is(target error) bool {
	_, ok := target.(*AlreadyExistsError)
	return ok
}
