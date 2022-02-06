package service

import (
	"errors"
	"net/http"

	"github.com/ashalfarhan/realworld/conduit"
)

var (
	// UserService Error
	ErrNoUserFound   = errors.New("no user found")
	ErrEmailExist    = errors.New("email already exist")
	ErrUsernameExist = errors.New("username already exist")
	ErrSelfFollow    = errors.New("you cannot follow your self")
	ErrSelfUnfollow  = errors.New("you cannot unfollow your self")
	ErrAlreadyFollow = errors.New("you are already follow this user")

	// AuthService Error
	ErrInvalidClaim = errors.New("invalid claim")

	// ArticleService Error
	ErrNoArticleFound          = errors.New("no article found")
	ErrNotAllowedDeleteArticle = errors.New("you cannot delete this article")
	ErrNotAllowedUpdateArticle = errors.New("you cannot edit this article")
	ErrNoCommentFound          = errors.New("no comment found")
	ErrNotAllowedDeleteComment = errors.New("you cannot delete this comment")
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
