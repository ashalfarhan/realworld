package middleware

import (
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/service"
	"github.com/sirupsen/logrus"
)

type ConduitMiddleware struct {
	authService *service.AuthService
	logger      *logrus.Entry
}

func NewMiddleware(s *service.Service) *ConduitMiddleware {
	return &ConduitMiddleware{
		s.AuthService,
		conduit.NewLogger("Service", "ConduitMiddleware"),
	}
}
