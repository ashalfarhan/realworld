package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	Addr string
	Port string
	// Should be "dev" or "prod"
	Env string
}

func LoadConfig() *Config {
	var ok bool
	co := &Config{}

	co.Port, ok = os.LookupEnv("PORT")
	if !ok {
		co.Port = "4000"
	}

	host := os.Getenv("HOST")

	co.Addr = fmt.Sprintf("%s:%s", host, co.Port)

	co.Env, ok = os.LookupEnv("APP_ENV")
	if !ok {
		co.Env = "dev"
		log.Println("APP_ENV is not set, using default \"dev\"")
	}

	return co
}
