package main

import (
	"github.com/ashalfarhan/realworld/api"
	"github.com/ashalfarhan/realworld/db"
)

func main() {
	conn := db.Connect()

	api.Bootstrap(conn)
}
