package main

import (
	"github.com/ashalfarhan/realworld/api"
	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/db"
	"github.com/jmoiron/sqlx"
)

var conn *sqlx.DB

func init() {
	config.Load()
	conn = db.Connect()
}

func main() {

	api.Bootstrap(conn)
}
