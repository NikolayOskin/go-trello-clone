package app

import (
	"github.com/NikolayOskin/go-trello-clone/controller"
	mid "github.com/NikolayOskin/go-trello-clone/controller/middleware"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/go-chi/chi"
)

func (a *App) InitRouting() {
	a.Router = chi.NewRouter()

	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.AllowContentType("application/json"))
	a.Router.Use(getCorsOpts().Handler)

	a.Router.Route("/auth", func(r chi.Router) {
		ctrl := &controller.AuthController{}
		r.Post("/sign-in", ctrl.SignIn)
		r.Post("/sign-up", ctrl.SignUp)
		r.Post("/reset-password", ctrl.ResetPassword)
		r.Post("/new-password", ctrl.NewPassword)
		r.With(mid.JWTCheck).Put("/verify/{code:[0-9]+}", ctrl.VerifyEmail)
	})

	a.Router.Route("/users", func(r chi.Router) {
		r.Use(mid.JWTCheck)
		ctrl := &controller.UserController{}
		r.Get("/me", ctrl.GetAuthUser)
		r.With(mid.Verified).Get("/me/boards", ctrl.GetBoards)
	})

	a.Router.Route("/boards", func(r chi.Router) {
		r.Use(mid.JWTCheck)
		r.Use(mid.Verified)
		ctrl := &controller.BoardController{}
		r.Get("/{id:[a-z0-9]+}", ctrl.GetFull)
		r.With(mid.DecodeBoardObj).Put("/{id:[a-z0-9]{24}}", ctrl.Update)
		r.Post("/", ctrl.Create)
	})

	a.Router.Route("/cards", func(r chi.Router) {
		r.Use(mid.JWTCheck)
		r.Use(mid.Verified)
		ctrl := &controller.CardController{}
		r.Post("/", ctrl.Create)
		r.With(mid.DecodeCardObj).Put("/{id:[a-z0-9]{24}}", ctrl.Update)
		r.Delete("/{id:[a-z0-9]{24}}", ctrl.Delete)
	})

	a.Router.Route("/lists", func(r chi.Router) {
		r.Use(mid.JWTCheck)
		r.Use(mid.Verified)
		ctrl := &controller.ListController{}
		r.Post("/", ctrl.Create)
		r.With(mid.DecodeListObj).Put("/{id:[a-z0-9]{24}}", ctrl.Update)
		r.Delete("/{id:[a-z0-9]{24}}", ctrl.Delete)
	})
}

func getCorsOpts() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
