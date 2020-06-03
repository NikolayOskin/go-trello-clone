package controller

import (
	"github.com/NikolayOskin/go-trello-clone/config"
	"github.com/NikolayOskin/go-trello-clone/handlers"
	"github.com/NikolayOskin/go-trello-clone/model"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

type AuthController struct{}

func (a *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := decodeJSONBody(w, r, &user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = handlers.HandleAuthenticate(&user)
	if err != nil {
		RespondJSON(w, 422, &Response{Message: err.Error()})
		return
	}

	token, err := a.generateJWTToken(user)
	if err != nil {
		RespondJSON(w, 422, &Response{Message: err.Error()})
		return
	}

	RespondJSON(w, 200, &JWTResponse{token})
}

func (a *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := decodeJSONBody(w, r, &user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = handlers.HandleAddUser(user)
	if err != nil {
		RespondJSON(w, 400, &Response{Message: err.Error()})
		return
	}

	RespondJSON(w, 200, &Response{Message: "Success"})
}

func (a *AuthController) generateJWTToken(user model.User) (string, error) {
	claims := model.JWTClaims{
		User: model.User{
			ID:    user.ID,
			Email: user.Email,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(72 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWTSecret))

	return tokenString, err
}
