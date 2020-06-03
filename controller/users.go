package controller

import (
	mid "github.com/NikolayOskin/go-trello-clone/middleware"
	"github.com/NikolayOskin/go-trello-clone/model"
	"net/http"
)

type UserController struct{}

func (u *UserController) GetAuthUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(mid.UserCtx).(model.User)
	RespondJSON(w, 200, user)
}
