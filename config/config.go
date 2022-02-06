package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	PgSource string
	Addr     string
	Port     string
	Env      string
}

var Co = &Config{}

func Load() {
	var ok bool

	Co.Port, ok = os.LookupEnv("PORT")
	if !ok {
		Co.Port = "4000"
	}
	Co.Env, ok = os.LookupEnv("APP_ENV")
	if !ok {
		Co.Env = "dev"
	}

	Co.Addr = fmt.Sprintf("%s:%s", os.Getenv("HOST"), Co.Port)
	Co.PgSource = fmt.
		Sprintf("user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"),
			os.Getenv("POSTGRES_SSLMODE"))
}
