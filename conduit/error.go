package conduit

import (
	"errors"
)

var (
	ErrInternal     = errors.New("internal server error")
	ErrUnauthorized = errors.New("unauthorized error")
	ErrForbidden    = errors.New("forbidden error")
	ErrNotFound     = errors.New("resource not found")

	// UserService Error
	ErrNoUserFound   = errors.New("no user found")
	ErrEmailExist    = errors.New("email already exist")
	ErrUsernameExist = errors.New("username already exist")
	ErrSelfFollow    = errors.New("you cannot follow your self")
	ErrSelfUnfollow  = errors.New("you cannot unfollow your self")
	ErrAlreadyFollow = errors.New("you are already follow this user")
)
