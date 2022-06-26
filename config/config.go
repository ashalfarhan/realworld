package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	PgSource  string
	RedisPass string
	Addr      string
	Port      string
	Env       string
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
	Co.PgSource = os.Getenv("POSTGRES_URL")
	Co.RedisPass = os.Getenv("REDIS_PASSWORD")
}
