package app

import (
	"github.com/NikolayOskin/go-trello-clone/controller"
	mid "github.com/NikolayOskin/go-trello-clone/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (a *App) InitRouting() {
	a.Router = chi.NewRouter()

	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.AllowContentType("application/json"))

	a.Router.Route("/auth", func(r chi.Router) {
		ctrl := &controller.AuthController{}
		r.Post("/sign-in", ctrl.SignIn)
		r.Post("/sign-up", ctrl.SignUp)
		r.With(mid.JWTCheck).Put("/verify/{code:[0-9]}", ctrl.VerifyEmail)
	})

	a.Router.Route("/users", func(r chi.Router) {
		r.Use(mid.JWTCheck)
		ctrl := &controller.UserController{}
		r.Get("/me", ctrl.GetAuthUser)
		r.Get("/me/boards", ctrl.GetBoards)
	})

	a.Router.Route("/boards", func(r chi.Router) {
		r.Use(mid.JWTCheck)
		ctrl := &controller.BoardController{}
		r.Post("/", ctrl.Create)
		r.With(mid.DecodeBoardObj).Put("/{id:[a-z0-9]{24}}", ctrl.Update)
	})

	a.Router.Route("/cards", func(r chi.Router) {
		r.Use(mid.JWTCheck)
		ctrl := &controller.CardController{}
		r.Post("/", ctrl.Create)
		r.With(mid.DecodeCardObj).Put("/{id:[a-z0-9]{24}}", ctrl.Update)
		r.Delete("/{id:[a-z0-9]{24}}", ctrl.Delete)
	})

	a.Router.Route("/lists", func(r chi.Router) {
		r.Use(mid.JWTCheck)
		ctrl := &controller.ListController{}
		r.Post("/", ctrl.Create)
		r.With(mid.DecodeListObj).Put("/{id:[a-z0-9]{24}}", ctrl.Update)
		r.Delete("/{id:[a-z0-9]{24}}", ctrl.Delete)
	})
}
