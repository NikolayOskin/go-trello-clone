package controller

import (
	"net/http"

	mid "github.com/NikolayOskin/go-trello-clone/controller/middleware"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/service/auth"
	"github.com/NikolayOskin/go-trello-clone/service/handlers"
	v "github.com/NikolayOskin/go-trello-clone/service/validator"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	JwtService  *auth.JWTService
	Validate    *validator.Validate
	AuthHandler handlers.Auth
	UserHandler handlers.User
}

func (a *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := decodeJSON(w, r, &user); err != nil {
		JSONResp(w, err.(*malformedRequest).Status, ErrResp{err.Error()})
		return
	}
	if user.Email == "" || user.Password == "" {
		JSONResp(w, 422, ErrResp{"username or password can't be empty"})
		return
	}
	if err := a.AuthHandler.Authenticate(r.Context(), &user); err != nil {
		JSONResp(w, 400, ErrResp{err.Error()})
		return
	}
	token, err := a.JwtService.GenerateToken(user)
	if err != nil {
		JSONResp(w, 400, ErrResp{err.Error()})
		return
	}
	JSONResp(w, 200, JWTResponse{token})
}

func (a *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := decodeJSON(w, r, &user); err != nil {
		JSONResp(w, err.(*malformedRequest).Status, ErrResp{err.Error()})
		return
	}
	if err := a.Validate.Struct(user); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			JSONResp(w, 422, ErrResp{e.Translate(v.Trans)})
			return
		}
	}
	if err := a.UserHandler.Create(r.Context(), user); err != nil {
		JSONResp(w, 400, ErrResp{err.Error()})
		return
	}
	JSONResp(w, 201, Response{Message: "Created"})
}

func (a *AuthController) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(mid.UserCtx).(model.User)
	code := chi.URLParam(r, "code")
	if err := a.AuthHandler.VerifyEmail(r.Context(), user, code); err != nil {
		JSONResp(w, 400, ErrResp{err.Error()})
		return
	}
	JSONResp(w, 200, Response{Message: "Verified"})
}

func (a *AuthController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := decodeJSON(w, r, &req); err != nil {
		JSONResp(w, err.(*malformedRequest).Status, ErrResp{err.Error()})
		return
	}
	if err := a.Validate.Struct(req); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			JSONResp(w, 422, ErrResp{e.Translate(v.Trans)})
			return
		}
	}
	if err := a.AuthHandler.ResetPassword(r.Context(), req.Email); err != nil {
		JSONResp(w, 400, ErrResp{err.Error()})
		return
	}
	JSONResp(w, 200, Response{Message: "Verification code sent"})
}

func (a *AuthController) NewPassword(w http.ResponseWriter, r *http.Request) {
	var req NewPasswordRequest
	if err := decodeJSON(w, r, &req); err != nil {
		JSONResp(w, err.(*malformedRequest).Status, ErrResp{err.Error()})
		return
	}
	if err := a.Validate.Struct(req); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			JSONResp(w, 422, ErrResp{e.Translate(v.Trans)})
			return
		}
	}
	if err := a.AuthHandler.SetNewPassword(r.Context(), req.Email, req.Code, req.Password); err != nil {
		JSONResp(w, 400, ErrResp{err.Error()})
		return
	}
	JSONResp(w, 200, Response{Message: "Password changed"})
}
