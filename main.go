package main

import (
	"github.com/NikolayOskin/go-trello-clone/app"
)

func main() {
	a := app.New()

	a.InitRouting()
	a.InitServices()
	a.Run()
}
