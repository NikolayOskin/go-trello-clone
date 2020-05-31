package app

import (
	"github.com/NikolayOskin/go-trello-clone/controller"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (a *App) InitRouting() {
	a.Router = chi.NewRouter()

	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.AllowContentType("application/json"))

	a.Router.Route("/auth", func(r chi.Router) {
		r.Post("/sign-in", controller.Authenticate)
	})

	a.Router.Route("/users", func(r chi.Router) {
		r.Get("/me", controller.GetAuthUser)
		r.Post("/", controller.AddUser)
	})

	a.Router.Route("/boards", func(r chi.Router) {
		r.Post("/", controller.AddBoard)
		r.Get("/", controller.GetBoards)
	})

	a.Router.Route("/cards", func(r chi.Router) {
		r.Post("/", controller.AddCard)
		r.Delete("/", controller.DeleteCard)
		r.Put("/", controller.UpdateCard)
	})

	//a.Router.Route("/lists", func(r chi.Router) {
	//	r.Post("/", controller.AddList)
	//})
	//
	//a.Router.Route("/cards", func(r chi.Router) {
	//	r.Post("/", controller.AddCard)
	//})
}

