package app

import (
	"os"

	"github.com/NikolayOskin/go-trello-clone/controller"
	mid "github.com/NikolayOskin/go-trello-clone/controller/middleware"
	"github.com/NikolayOskin/go-trello-clone/service/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/go-chi/chi"
)

func (app *app) InitRouting() {
	const idPattern = "/{id:[a-z0-9]{24}}"

	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.AllowContentType("application/json"))
	a.Router.Use(getCorsOpts().Handler)

	app.Router.Use(middleware.AllowContentType("application/json"))
	app.Router.Use(getCorsOpts().Handler)

	app.Router.Route("/auth", func(r chi.Router) {
		ctrl := controller.NewAuthCtrl(app.JWTService, app.Validator, handlers.Auth{}, handlers.User{})

		r.Post("/sign-in", ctrl.SignIn)
		r.Post("/sign-up", ctrl.SignUp)
		r.Post("/reset-password", ctrl.ResetPassword)
		r.Post("/new-password", ctrl.NewPassword)
		r.With(mid.JWTCheck(app.JWTService)).Put("/verify/{code:[0-9]+}", ctrl.VerifyEmail)
	})

	app.Router.Route("/users", func(r chi.Router) {
		ctrl := controller.NewUserCtrl()

		r.Use(mid.JWTCheck(app.JWTService))
		r.Get("/me", ctrl.GetAuthUser)
		r.With(mid.Verified).Get("/me/boards", ctrl.GetBoards)
	})

	app.Router.Route("/boards", func(r chi.Router) {
		ctrl := controller.NewBoardCtrl(app.Validator, handlers.Board{})

		r.Use(mid.JWTCheck(app.JWTService), mid.Verified)
		r.Get(idPattern, ctrl.GetFull)
		r.With(mid.DecodeBoardObj(app.Validator)).Put(idPattern, ctrl.Update)
		r.Post("/", ctrl.Create)
	})

	app.Router.Route("/cards", func(r chi.Router) {
		ctrl := controller.NewCardCtrl()

		r.Use(mid.JWTCheck(app.JWTService), mid.Verified)
		r.Post("/", ctrl.Create)
		r.With(mid.DecodeCardObj(app.Validator)).Put(idPattern, ctrl.Update)
		r.Delete(idPattern, ctrl.Delete)
	})

	app.Router.Route("/lists", func(r chi.Router) {
		ctrl := controller.NewListCtrl(app.Validator)

		r.Use(mid.JWTCheck(app.JWTService), mid.Verified)
		r.Post("/", ctrl.Create)
		r.With(mid.DecodeListObj(app.Validator)).Put(idPattern, ctrl.Update)
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
