package main

import (
	"github.com/NikolayOskin/go-trello-clone/app"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"log"
	"net/http"
)

func main() {
	a := app.App{}
	mongodb.InitDB()
	a.InitRouting()

	err := http.ListenAndServe(":3001", a.Router)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}
