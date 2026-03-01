package main

import (
	"to-do-list/internal/repository/db"
	"to-do-list/internal/routes"
	"to-do-list/pkg/conf"
)

func main() {
	conf.Init()
	db.InitDB()

	r := routes.NewRouter()

	r.Run(":3000")
}
