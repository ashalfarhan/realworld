package persistence

import (
	"github.com/ashalfarhan/realworld/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func Connect() *sqlx.DB {
	conn, err := sqlx.Open("postgres", config.PgSource)
	if err != nil {
		logrus.Panicln("Failed to opening database connection:", err)
	}

	if err = conn.Ping(); err != nil {
		defer conn.Close()
		logrus.Panicln("Failed to connect to database:", err)
	}

	logrus.Println("Successfully connected to database")
	driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil {
		logrus.Panicln("Failed to initialize postgres migration:", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+config.MigrationPath, "postgres", driver)
	if err != nil {
		logrus.Panicln("Failed to initialize migration instance:", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Panicln("Failed to run migration", err)
	}
	logrus.Println("Successfully run all migrations")
	return conn
}
