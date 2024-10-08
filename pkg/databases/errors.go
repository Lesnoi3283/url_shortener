package databases

import "errors"

// лучше так:
//var ErrAlreadyExists error = errors.New("already exists")
//func NewAlreadyExistsError()

// AlreadyExistsError is a custom error witch can return a ShortURL of URL you tried to save.
type AlreadyExistsError struct {
	ShortURL string
	Text     string
}

// Error returns a ShortURL.
func (a *AlreadyExistsError) Error() string {
	return a.ShortURL
}

// NewAlreadyExistsError creates a new AlreadyExistsError.
func NewAlreadyExistsError(shortURL string) *AlreadyExistsError {
	return &AlreadyExistsError{ShortURL: shortURL}
}

// Is is same as errors.Is func but for an AlreadyExistsError.
func (a *AlreadyExistsError) Is(target error) bool {
	_, ok := target.(*AlreadyExistsError)
	return ok
}

var errURLWasDeleted = errors.New("this url was marked as deleted")

// ErrURLWasDeleted returns an errURLWasDeleted error.
func ErrURLWasDeleted() error {
	return errURLWasDeleted
}

var errThisFuncIsNotSupported = errors.New("jsonFileStorage storage doesnt support deleteBatch func yet")

// ErrThisFuncIsNotSupported returns an errThisFuncIsNotSupported error.
func ErrThisFuncIsNotSupported() error {
	return errThisFuncIsNotSupported
}
