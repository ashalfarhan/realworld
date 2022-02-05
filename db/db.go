package db

import (
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect() *sqlx.DB {
	conn, err := sqlx.Open("postgres", "user=postgres password=password dbname=postgres sslmode=disable")
	if err != nil {
		conduit.Logger.Panicf("Failed to opening database connection: %v\n", err)
		return nil
	}

	if err = conn.Ping(); err != nil {
		defer conn.Close()
		conduit.Logger.Panicf("Failed to connect to database %v\n", err)
		return nil
	}

	conduit.Logger.Println("Successfully connected to database")
	return conn
}
