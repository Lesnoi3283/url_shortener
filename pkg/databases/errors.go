package databases

import "errors"

// лучше так:
//var ErrAlreadyExists error = errors.New("already exists")

//func NewAlreadyExistsError()

type AlreadyExistsError struct {
	ShortURL string
	Text     string
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

var errURLWasDeleted = errors.New("this url was marked as deleted")

func ErrURLWasDeleted() error {
	return errURLWasDeleted
}

var errThisFuncIsNotSupported = errors.New("jsonFileStorage storage doesnt support deleteBatch func yet")

func ErrThisFuncIsNotSupported() error {
	return errThisFuncIsNotSupported
}
