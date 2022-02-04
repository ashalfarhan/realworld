package db

import (
	"database/sql"

	"github.com/ashalfarhan/realworld/conduit"
	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	connection, err := sql.Open("postgres", "user=postgres password=password dbname=postgres sslmode=disable")
	if err != nil {
		conduit.Logger.Panicf("Failed to opening database connection: %v\n", err)
		return nil
	}

	if err = connection.Ping(); err != nil {
		defer connection.Close()
		conduit.Logger.Panicf("Failed to connect to database %v\n", err)
		return nil
	}

	conduit.Logger.Println("Successfully connected to database")
	return connection
}
