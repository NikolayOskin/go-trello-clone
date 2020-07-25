package main

import (
	"github.com/NikolayOskin/go-trello-clone/app"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"log"
)

func main() {
	a := app.New()
	mongodb.InitDB()

	a.InitRouting()
	a.InitServices()

	log.Println("Starting server...")

	a.Run(":3001")
}
