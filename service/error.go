package service

import (
	"net/http"

	"github.com/ashalfarhan/realworld/conduit"
)

const (
	ErrDuplicateEmail    = "pq: duplicate key value violates unique constraint \"users_email_key\""
	ErrDuplicateUsername = "pq: duplicate key value violates unique constraint \"users_username_key\""
	ErrDuplicateFollowing = "pq: duplicate key value violates unique constraint \"followings_pkey\""
)

type ServiceError struct {
	Code  int
	Error error
}

func CreateServiceError(code int, original error) *ServiceError {
	if code == http.StatusInternalServerError {
		return &ServiceError{code, conduit.ErrInternal}
	}

	return &ServiceError{code, original}
}
