package config

import (
	"fmt"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	PgSource  string
	RedisPass string
	Addr      string
	Port      string
	Env       string
	CacheTTL  time.Duration
}

var Co = &Config{}

func Load() {
	var ok bool

	if Co.Port, ok = os.LookupEnv("PORT"); !ok {
		Co.Port = "4000"
	}

	if Co.Env, ok = os.LookupEnv("APP_ENV"); !ok {
		Co.Env = "dev"
	}

	Co.CacheTTL = time.Microsecond * 2 // Change to microsecond if testing with postman spec
	Co.Addr = fmt.Sprintf("%s:%s", os.Getenv("HOST"), Co.Port)
	Co.PgSource = os.Getenv("POSTGRES_URL")
	Co.RedisPass = os.Getenv("REDIS_PASSWORD")
}
