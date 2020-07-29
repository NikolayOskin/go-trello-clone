package app

import (
	"github.com/NikolayOskin/go-trello-clone/controller"
	mid "github.com/NikolayOskin/go-trello-clone/controller/middleware"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/go-chi/chi"
)

func (a *app) InitRouting() {
	const idPattern = "/{id:[a-z0-9]{24}}"

	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.AllowContentType("application/json"))
	a.Router.Use(getCorsOpts().Handler)

	a.Router.Route("/auth", func(r chi.Router) {
		ctrl := &controller.AuthController{AuthService: a.Auth, Validate: a.Validator}

		r.Post("/sign-in", ctrl.SignIn)
		r.Post("/sign-up", ctrl.SignUp)
		r.Post("/reset-password", ctrl.ResetPassword)
		r.Post("/new-password", ctrl.NewPassword)
		r.With(mid.JWTCheck(a.Auth)).Put("/verify/{code:[0-9]+}", ctrl.VerifyEmail)
	})

	a.Router.Route("/users", func(r chi.Router) {
		ctrl := &controller.UserController{}

		r.Use(mid.JWTCheck(a.Auth))
		r.Get("/me", ctrl.GetAuthUser)
		r.With(mid.Verified).Get("/me/boards", ctrl.GetBoards)
	})

	a.Router.Route("/boards", func(r chi.Router) {
		ctrl := &controller.BoardController{Validate: a.Validator}

		r.Use(mid.JWTCheck(a.Auth), mid.Verified)
		r.Get("/{id:[a-z0-9]+}", ctrl.GetFull)
		r.With(mid.DecodeBoardObj(a.Validator)).Put(idPattern, ctrl.Update)
		r.Post("/", ctrl.Create)
	})

	a.Router.Route("/cards", func(r chi.Router) {
		ctrl := &controller.CardController{}

		r.Use(mid.JWTCheck(a.Auth), mid.Verified)
		r.Post("/", ctrl.Create)
		r.With(mid.DecodeCardObj(a.Validator)).Put(idPattern, ctrl.Update)
		r.Delete(idPattern, ctrl.Delete)
	})

	a.Router.Route("/lists", func(r chi.Router) {
		ctrl := &controller.ListController{Validate: a.Validator}

		r.Use(mid.JWTCheck(a.Auth), mid.Verified)
		r.Post("/", ctrl.Create)
		r.With(mid.DecodeListObj(a.Validator)).Put(idPattern, ctrl.Update)
		r.Delete(idPattern, ctrl.Delete)
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
