package main

import (
	"github.com/ashalfarhan/realworld/api"
	"github.com/ashalfarhan/realworld/cache"
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/db"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var conn *sqlx.DB

func init() {
	config.Load()
	conduit.Logger = logrus.WithField("context", "ConduitApp")
	cache.Init()
}

func main() {
	conn = db.Connect()
	api.Bootstrap(conn)
}
