package api

import (
	"net/http"
	"time"

	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/service"
)

func InitServer(serv *service.Service) *http.Server {
	r := InitRoutes(serv)
	return &http.Server{
		Addr:         config.Addr,
		Handler:      r,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
}
