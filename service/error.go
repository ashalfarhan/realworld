package service

import (
	"errors"
)

var (
	// UserService Error
	ErrNoUserFound   = errors.New("no user found")
	ErrEmailExist    = errors.New("email already exist")
	ErrUsernameExist = errors.New("username already exist")
	ErrIdentityExist = errors.New("username or email is in use")
	ErrSelfFollow    = errors.New("you cannot follow your self")
	ErrSelfUnfollow  = errors.New("you cannot unfollow your self")
	ErrAlreadyFollow = errors.New("you are already follow this user")

	// AuthService Error
	ErrInvalidClaim    = errors.New("invalid claim")
	ErrInvalidIdentity = errors.New("invalid identity or password")

	// ArticleService Error
	ErrNoArticleFound          = errors.New("no article found")
	ErrNotAllowedDeleteArticle = errors.New("you cannot delete this article")
	ErrNotAllowedUpdateArticle = errors.New("you cannot edit this article")
	ErrNoCommentFound          = errors.New("no comment found")
	ErrNotAllowedDeleteComment = errors.New("you cannot delete this comment")
)
