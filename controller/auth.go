package controller

import (
	"github.com/NikolayOskin/go-trello-clone/handlers"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"
)

type AuthController struct{}

func (a *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := decodeJSON(w, r, &user); err != nil {
		JSONResp(w, 400, &ErrResp{err.Error()})
		return
	}
	if err := handlers.Authenticate(&user); err != nil {
		JSONResp(w, 422, &ErrResp{err.Error()})
		return
	}
	token, err := a.generateJWTToken(user)
	if err != nil {
		JSONResp(w, 422, &ErrResp{err.Error()})
		return
	}
	JSONResp(w, 200, &JWTResponse{token})
}

func (a *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := decodeJSON(w, r, &user); err != nil {
		JSONResp(w, 400, &ErrResp{err.Error()})
		return
	}
	if err := handlers.CreateUser(user); err != nil {
		JSONResp(w, 400, &ErrResp{err.Error()})
		return
	}
	JSONResp(w, 200, &Response{Message: "Success"})
}

func (a *AuthController) VerifyEmail(w http.ResponseWriter, r *http.Request) {

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
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return tokenStr, err
}
