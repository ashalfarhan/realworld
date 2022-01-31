package db

import (
	"database/sql"

	"github.com/ashalfarhan/realworld/conduit"
	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	connection, err := sql.Open("postgres", "user=postgres password=password dbname=postgres sslmode=disable")
	if err != nil {
		conduit.Logger.Panicf("failed to opening database connection: %s\n", err.Error())
		return nil
	}

	if err = connection.Ping(); err != nil {
		defer connection.Close()
		conduit.Logger.Panicf("failed to connect to database %s\n", err.Error())
		return nil
	}

	conduit.Logger.Println("Successfully connected to database")
	return connection
}
