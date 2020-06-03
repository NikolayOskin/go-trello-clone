package controller

import (
	"github.com/NikolayOskin/go-trello-clone/handlers"
	mid "github.com/NikolayOskin/go-trello-clone/middleware"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/repository"
	"net/http"
)

type BoardController struct{}

func (b *BoardController) AddBoard(w http.ResponseWriter, r *http.Request) {
	var board model.Board
	user := r.Context().Value(mid.UserCtx).(model.User)

	err := decodeJSONBody(w, r, &board)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	board.UserId = user.ID.Hex()
	handlers.HandleAddBoard(board)

	RespondJSON(w, 200, &Response{Message: "Success"})
}

func (b *BoardController) GetBoards(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(mid.UserCtx).(model.User)
	boardsRepo := repository.Boards{}

	RespondJSON(w, 200, boardsRepo.FetchByUser(user))
}
