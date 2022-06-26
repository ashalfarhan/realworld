package main

import (
	"github.com/ashalfarhan/realworld/api"
	"github.com/ashalfarhan/realworld/cache"
	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/persistence"
	"github.com/ashalfarhan/realworld/utils/logger"
	"github.com/jmoiron/sqlx"
)

var conn *sqlx.DB

func init() {
	config.Load()
	logger.Init()
	cache.Init()
}

func main() {
	conn = persistence.Connect()
	api.Bootstrap(conn)
}
