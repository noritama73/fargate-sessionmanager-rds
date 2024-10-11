package service

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrInvalid    = errors.New("invalid")
	ErrUnexpected = errors.New("unexpected error")
)

type ServiceError struct {
	Code          error
	OriginalError error
}

func (e ServiceError) Error() string {
	return e.OriginalError.Error()
}

func (e ServiceError) ErrCode() error {
	return e.Code
}
