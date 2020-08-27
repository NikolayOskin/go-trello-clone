package controller

import (
	"github.com/NikolayOskin/go-trello-clone/service/auth"
	"github.com/NikolayOskin/go-trello-clone/service/handlers"
	"github.com/go-playground/validator/v10"
)

func NewAuthCtrl(
	jwtService *auth.JWTService,
	validator *validator.Validate,
	authHandler handlers.Auth,
	userHandler handlers.User,
) *AuthController {
	return &AuthController{
		JwtService:  jwtService,
		Validate:    validator,
		AuthHandler: authHandler,
		UserHandler: userHandler,
	}
}

func NewBoardCtrl(
	validator *validator.Validate,
	boardHandler handlers.Board,
) *BoardController {
	return &BoardController{
		Validate:     validator,
		BoardHandler: boardHandler,
	}
}

func NewUserCtrl() *UserController {
	return &UserController{}
}

func NewCardCtrl() *CardController {
	return &CardController{}
}

func NewListCtrl(validator *validator.Validate) *ListController {
	return &ListController{
		Validate: validator,
	}
}
