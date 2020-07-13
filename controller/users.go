package controller

import (
	mid "github.com/NikolayOskin/go-trello-clone/controller/middleware"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/repository"
	"net/http"
)

type UserController struct{}

func (u *UserController) GetAuthUser(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(mid.UserCtx).(model.User)
	repo := repository.Users{}
	user, err := repo.GetById(userCtx.ID.Hex())
	if err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Could not fetch user"})
	}
	if user != nil {
		userCtx.Verified = user.Verified
	}

	JSONResp(w, 200, userCtx)
}

func (u *UserController) GetBoards(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(mid.UserCtx).(model.User)
	boardsRepo := repository.Boards{}
	boards, err := boardsRepo.FetchByUser(user)
	if err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Could not fetch user boards"})
	}
	JSONResp(w, 200, boards)
}
