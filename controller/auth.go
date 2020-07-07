package controller

import (
	mid "github.com/NikolayOskin/go-trello-clone/controller/middleware"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/service/handlers"
	v "github.com/NikolayOskin/go-trello-clone/service/validator"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"net/http"
	"os"
	"time"
)

type AuthController struct{}

func (a *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := decodeJSON(w, r, &user); err != nil {
		JSONResp(w, err.(*malformedRequest).Status, &ErrResp{err.Error()})
		return
	}
	if user.Email == "" || user.Password == "" {
		JSONResp(w, 422, &ErrResp{"username or password can't be empty"})
		return
	}
	if err := handlers.Authenticate(&user, r.Context()); err != nil {
		JSONResp(w, 400, &ErrResp{err.Error()})
		return
	}
	token, err := a.generateJWTToken(user)
	if err != nil {
		JSONResp(w, 400, &ErrResp{err.Error()})
		return
	}
	JSONResp(w, 200, &JWTResponse{token})
}

func (a *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := decodeJSON(w, r, &user); err != nil {
		JSONResp(w, err.(*malformedRequest).Status, &ErrResp{err.Error()})
		return
	}
	validate := v.New()
	if err := validate.Struct(user); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			JSONResp(w, 422, &ErrResp{e.Translate(v.Trans)})
			return
		}
	}
	if err := handlers.CreateUser(user, r.Context()); err != nil {
		JSONResp(w, 400, &ErrResp{err.Error()})
		return
	}
	JSONResp(w, 201, &Response{Message: "Created"})
}

func (a *AuthController) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(mid.UserCtx).(model.User)
	code := chi.URLParam(r, "code")
	if code == "" {
		JSONResp(w, 422, &ErrResp{"code must not be empty"})
		return
	}
	if err := handlers.VerifyEmail(user, code, r.Context()); err != nil {
		JSONResp(w, 400, &ErrResp{err.Error()})
		return
	}
	JSONResp(w, 200, &Response{Message: "Verified"})
}

func (a *AuthController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := decodeJSON(w, r, &req); err != nil {
		JSONResp(w, err.(*malformedRequest).Status, &ErrResp{err.Error()})
		return
	}
	validate := v.New()
	if err := validate.Struct(req); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			JSONResp(w, 422, &ErrResp{e.Translate(v.Trans)})
			return
		}
	}
	if err := handlers.ResetPassword(req.Email, r.Context()); err != nil {
		JSONResp(w, 400, &ErrResp{err.Error()})
		return
	}
	JSONResp(w, 200, &Response{Message: "Verification code sent"})
}

func (a *AuthController) NewPassword(w http.ResponseWriter, r *http.Request) {
	var req NewPasswordRequest
	if err := decodeJSON(w, r, &req); err != nil {
		JSONResp(w, err.(*malformedRequest).Status, &ErrResp{err.Error()})
		return
	}
	validate := v.New()
	if err := validate.Struct(req); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			JSONResp(w, 422, &ErrResp{e.Translate(v.Trans)})
			return
		}
	}
	if err := handlers.SetNewPassword(req.Email, req.Code, req.Password); err != nil {
		JSONResp(w, 400, &ErrResp{err.Error()})
		return
	}
	JSONResp(w, 200, &Response{Message: "Password changed"})
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
