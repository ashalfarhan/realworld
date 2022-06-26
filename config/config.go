package config

import (
	"fmt"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var (
	PgSource  string
	RedisPass string
	Addr      string
	Port      string
	Env       string
	CacheTTL  time.Duration
)

func Load() {
	var ok bool

	if Port, ok = os.LookupEnv("PORT"); !ok {
		Port = "4000"
	}

	if Env, ok = os.LookupEnv("APP_ENV"); !ok {
		Env = "dev"
	}

	CacheTTL = time.Millisecond * 2 // Change to microsecond if testing with postman spec
	Addr = fmt.Sprintf("%s:%s", os.Getenv("HOST"), Port)
	PgSource = os.Getenv("POSTGRES_URL")
	RedisPass = os.Getenv("REDIS_PASSWORD")
}
