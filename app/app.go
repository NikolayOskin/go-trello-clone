package app

import (
	"github.com/go-chi/chi"
	"net/http"
)

type App struct {
	Router *chi.Mux
}

func (a *App) Run(addr string) {
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		panic(err)
	}
}