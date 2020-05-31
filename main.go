package main

import (
	"github.com/NikolayOskin/go-trello-clone/app"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"net/http"
)

func main() {
	a := app.App{}

	mongodb.InitDB()

	a.InitRouting()

	err := http.ListenAndServe(":3000", a.Router)
	if err != nil {
		panic(err)
	}
}