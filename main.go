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
)

func init() {
	config.Load()
	logger.Init()
}

func main() {
	db := persistence.Connect()
	cacheStore := cache.Init()

	services := service.InitService(db, cacheStore)
	server := api.InitServer(services)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	logger.Log.Println("Booting up the server...")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Log.Errorf("Failed to start the server: %v", err)
			if err = cleanup(cacheStore, db); err != nil {
				logger.Log.Errorf("Failed to cleanup: %v", err)
			}
		}
	}()

	logger.Log.Printf("Listening on %s in %q mode", config.Co.Addr, config.Co.Env)

	<-shutdown
	logger.Log.Println("Gracefully shutdown...")

	defer func() {
		if err := cleanup(cacheStore, db); err != nil {
			logger.Log.Errorf("Failed to cleanup: %v", err)
		}
	}()

	if err := server.Close(); err != nil {
		logger.Log.Errorf("Failed to close the server: %v", err)
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
