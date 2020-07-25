package main

import (
	"github.com/NikolayOskin/go-trello-clone/app"
	"github.com/NikolayOskin/go-trello-clone/db"
)

func main() {
	a := app.New()
	db.InitDB()

	a.InitRouting()
	a.InitServices()

	log.Println("Starting server...")

	a.Run(":3001")
}
