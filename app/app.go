package app

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NikolayOskin/go-trello-clone/db"
	mailer "github.com/NikolayOskin/go-trello-clone/service"
	"github.com/NikolayOskin/go-trello-clone/service/auth"
	"github.com/NikolayOskin/go-trello-clone/service/validator"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	v "github.com/go-playground/validator/v10"
)

type app struct {
	Router     *chi.Mux
	JWTService *auth.JWTService
	Validator  *v.Validate
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

	// Initiating JWTService service
	const jwtTTLHours = 72 // jwt token life time in hours
	JWTService, err := auth.NewJWTService(privateKey, publicKey, jwtTTLHours)
	if err != nil {
		log.Fatal(err)
	}

	return &app{
		Router:     chi.NewRouter(),
		JWTService: JWTService,
		Validator:  validator.New(),
	}
}

func (a *app) Run() {
	serverPort := os.Getenv("HTTP_SERVER_PORT")
	if serverPort == "" {
		log.Fatal("HTTP_SERVER_PORT env is not set")
	}

	db.InitDB()
	defer db.Disconnect()

	log.Println("Ready to start server...")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	srv := http.Server{
		Addr:    net.JoinHostPort("", serverPort),
		Handler: a.Router,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("Server started. API listening on %s", net.JoinHostPort("", serverPort))
	}()

	select {
	case sig := <-shutdown:
		log.Printf("%v : Shutting down the server...", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		_ = srv.Close()
		log.Println("Could not stop the server gracefully")
		return
	}
	log.Println("Server stopped")
}

func (a *app) InitServices() {
	mailer.Start()
}
