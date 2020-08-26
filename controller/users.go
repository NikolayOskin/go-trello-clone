package controller

import (
	"net/http"

	mid "github.com/NikolayOskin/go-trello-clone/controller/middleware"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/repository"
)

type UserController struct{}

func (u *UserController) GetAuthUser(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(mid.UserCtx).(model.User)
	user, err := repository.Users.FindById(r.Context(), userCtx.ID.Hex())
	if err != nil {
		JSONResp(w, 500, ErrResp{Message: "Could not fetch user"})
		return
	}
	readUser := model.ReadUser{
		ID:    userCtx.ID,
		Email: userCtx.Email,
	}
	if user != nil {
		readUser.Verified = user.Verified
	}

	JSONResp(w, 200, readUser)
}

func (u *UserController) GetBoards(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(mid.UserCtx).(model.User)
	boards, err := repository.Boards.FetchByUser(r.Context(), user)
	if err != nil {
		JSONResp(w, 500, ErrResp{Message: "Could not fetch user boards"})
		return
	}
	JSONResp(w, 200, boards)
}
