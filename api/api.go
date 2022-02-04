package api

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/service"
	"github.com/gorilla/handlers"
)

func Bootstrap(db *sql.DB) {
	cf := config.LoadConfig()
	serv := service.InitService(db)
	server := InitServer(serv, cf)

	conduit.Logger.Printf("Listening on %s in \"%s\" mode", cf.Addr, cf.Env)

	if err := server.ListenAndServe(); err != nil {
		defer db.Close()
		conduit.Logger.Panicf("Failed to bootstrap the server: %v", err)
		return
	}
}

func InitServer(serv *service.Service, cf *config.Config) *http.Server {
	r := InitRoutes(serv)

	out := os.Stdout
	var h http.Handler

	if cf.Env == "dev" {
		h = handlers.LoggingHandler(out, r)
	} else {
		h = handlers.CombinedLoggingHandler(out, r)
	}

	return &http.Server{
		Addr:         cf.Addr,
		Handler:      h,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
}
