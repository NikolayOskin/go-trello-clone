package app

import (
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"log"
	"net/http"

	mailer "github.com/NikolayOskin/go-trello-clone/service"
	"github.com/NikolayOskin/go-trello-clone/service/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
)

type app struct {
	Router    *chi.Mux
	Auth      *auth.Auth
	Validator *validator.Validate
}

func New() *app {
	log.Println("Reading private.pem & public.pem files")

	privatePEM, err := ioutil.ReadFile("./private.pem")
	if err != nil {
		log.Fatal("cannot read private key from file")
	}
	publicPEM, err := ioutil.ReadFile("./public.pem")
	if err != nil {
		log.Fatal("cannot read public key from file")
	}

	// Parsing private & public keys
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		log.Fatal("cannot parse private key")
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		log.Fatal("cannot parse public key")
	}

	// Initiating Auth service
	const jwtTTLHours = 72 // jwt token life time in hours
	a, err := auth.New(privateKey, publicKey, jwtTTLHours)
	if err != nil {
		log.Fatal(err)
	}

	return &app{
		Router:    chi.NewRouter(),
		Auth:      a,
		Validator: validator.New(),
	}
}

func (a *app) RunServer(addr string) {
	log.Println("Starting server...")

	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}

func (a *app) InitServices() {
	mailer.Start()
}
