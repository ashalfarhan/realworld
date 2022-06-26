package persistence

import (
	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/utils/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect() *sqlx.DB {
	conn, err := sqlx.Open("postgres", config.PgSource)
	if err != nil {
		logger.Log.Panicf("Failed to opening database connection: %v\n", err)
		return nil
	}

	if err = conn.Ping(); err != nil {
		defer conn.Close()
		logger.Log.Panicf("Failed to connect to database %v\n", err)
		return nil
	}

	logger.Log.Println("Successfully connected to database")
	return conn
}
