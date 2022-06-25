package middleware

import (
	"github.com/ashalfarhan/realworld/service"
)

type ConduitMiddleware struct {
	authService *service.AuthService
}

func NewMiddleware(s *service.Service) *ConduitMiddleware {
	return &ConduitMiddleware{
		authService: s.AuthService,
	}
}
