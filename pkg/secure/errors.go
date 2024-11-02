package secure

import "errors"

var errTokenIsNotValid = errors.New("token is not valid")

func NewErrTokenIsNotValid() error {
	return errTokenIsNotValid
}
