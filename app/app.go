package app

import (
	mailer "github.com/NikolayOskin/go-trello-clone/service"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

type App struct {
	Router *chi.Mux
}

func (a *App) Run(addr string) {
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}

func (a *App) InitServices() {
	mailer.Start()
}
