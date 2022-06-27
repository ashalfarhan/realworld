package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ashalfarhan/realworld/api"
	"github.com/ashalfarhan/realworld/cache"
	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/persistence"
	"github.com/ashalfarhan/realworld/service"
	"github.com/ashalfarhan/realworld/utils/logger"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func init() {
	config.Load()
	logger.Configure()
}

func main() {
	db := persistence.Connect()
	store := cache.Init()

	services := service.InitService(db, store)
	server := api.InitServer(services)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	logrus.Println("Booting up the server...")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.Errorf("Failed to start the server: %v", err)
			if err = cleanup(store, db); err != nil {
				logrus.Errorf("Failed to cleanup: %v", err)
			}
		}
	}()

	logrus.Printf("Listening on %s in %q mode", config.Addr, config.Env)

	<-shutdown
	logrus.Println("Gracefully shutdown...")

	defer func() {
		if err := cleanup(store, db); err != nil {
			logrus.Errorf("Failed to cleanup: %v", err)
		}
	}()

	if err := server.Close(); err != nil {
		logrus.Errorf("Failed to close the server: %v", err)
	}
}

func cleanup(store *redis.Client, db *sqlx.DB) error {
	if err := store.Close(); err != nil {
		return fmt.Errorf("failed to close redis: %w", err)
	}
	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close postgres: %w", err)
	}
	return nil
}
