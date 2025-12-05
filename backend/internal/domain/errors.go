package domain

import "errors"

var ErrNotFound = errors.New("not found")

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(msg string) error {
	return &ValidationError{Message: msg}
}

func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}
