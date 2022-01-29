package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Bootstrap(db *sql.DB) {
	cf := config.LoadConfig()
	serv := service.InitService(db)
	r := InitRoutes(serv)
	s := InitServer(r, cf)

	log.Printf("Listening on %d\n", 4000)

	if err := s.ListenAndServe(); err != nil {
		defer db.Close()
		log.Fatalf("Failed to bootstrap the server: %s\n", err.Error())
		os.Exit(1)
	}
}

func InitServer(r *mux.Router, cf *config.Config) *http.Server {
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
