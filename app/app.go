package app

import (
	mailer "github.com/NikolayOskin/go-trello-clone/service"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

type app struct {
	Router *chi.Mux
}

func New() *app {
	return &app{}
}

func (a *app) Run(addr string) {
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}

func (a *app) InitServices() {
	mailer.Start()
}
