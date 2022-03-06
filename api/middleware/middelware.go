package middleware

import (
	"log"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/service"
)

type ConduitMiddleware struct {
	authService *service.AuthService
	logger *log.Logger
}

func NewMiddleware(s *service.Service) *ConduitMiddleware {
	return &ConduitMiddleware{
		s.AuthService,
		conduit.NewLogger("ConduitMiddleware"),
	}
}
