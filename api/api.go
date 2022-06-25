package api

import (
	"net/http"
	"time"

	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/service"
	"github.com/ashalfarhan/realworld/utils/logger"
	"github.com/jmoiron/sqlx"
)

func Bootstrap(db *sqlx.DB) {
	services := service.InitService(db)
	server := InitServer(services)

	logger.Log.Printf("Listening on %s in \"%s\" mode", config.Co.Addr, config.Co.Env)
	if err := server.ListenAndServe(); err != nil {
		defer db.Close()
		logger.Log.Panicf("Failed to bootstrap the server: %v", err)
		return
	}
}

func InitServer(serv *service.Service) *http.Server {
	r := InitRoutes(serv)

	return &http.Server{
		Addr:         config.Co.Addr,
		Handler:      r,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
}
