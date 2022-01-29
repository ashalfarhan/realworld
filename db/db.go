package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	connection, err := sql.Open("postgres", "user=postgres password=password dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalf("failed to opening database connection: %s\n", err.Error())
		os.Exit(1)
		return nil
	}

	if err = connection.Ping(); err != nil {
		log.Fatalf("failed to connect to database %s\n", err.Error())
		connection.Close()
		os.Exit(1)
		return nil
	}

	log.Println("Successfully connected to database")
	return connection
}
