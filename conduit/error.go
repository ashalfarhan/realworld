package conduit

import (
	"errors"

	"github.com/ashalfarhan/realworld/model"
)

var (
	ErrInternal     = errors.New("internal server error")
	ErrUnauthorized = errors.New("unauthorized error")
	ErrForbidden    = errors.New("forbidden error")
	ErrNotFound     = errors.New("resource not found")
	GeneralError    = BuildError(500, ErrInternal)
)

func BuildError(code int, original error) *model.ConduitError {
	return &model.ConduitError{
		Code: code,
		Err:  original,
	}
}
