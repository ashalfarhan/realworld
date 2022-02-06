package conduit

import (
	"errors"
)

var (
	ErrInternal     = errors.New("internal server error")
	ErrUnauthorized = errors.New("unauthorized error")
	ErrForbidden    = errors.New("forbidden error")
	ErrNotFound     = errors.New("resource not found")
)
