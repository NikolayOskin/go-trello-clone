package controller

import (
	"github.com/NikolayOskin/go-trello-clone/handlers"
	mid "github.com/NikolayOskin/go-trello-clone/middleware"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type BoardController struct{}

func (b *BoardController) GetById(w http.ResponseWriter, r *http.Request) {

}

func (b *BoardController) Create(w http.ResponseWriter, r *http.Request) {
	var board model.Board
	user := r.Context().Value(mid.UserCtx).(model.User)

	err := decodeJSONBody(w, r, &board)
	if err != nil {
		RespondJSON(w, 500, &ErrResp{Message: "Could not parse json request"})
		return
	}
	board.UserId = user.ID.Hex()
	err = handlers.HandleCreateBoard(board)
	if err != nil {
		RespondJSON(w, 500, &ErrResp{Message: "Server error"})
		return
	}

	RespondJSON(w, 200, &Response{Message: "Added"})
}

func (b *BoardController) Update(w http.ResponseWriter, r *http.Request) {
	board := r.Context().Value(mid.BoardCtx).(model.Board)

	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		RespondJSON(w, 500, &ErrResp{Message: "Server error"})
		return
	}
	board.ID = id
	err = handlers.HandleUpdateBoard(board)
	if err != nil {
		RespondJSON(w, 200, &ErrResp{Message: "Server error"})
		return
	}

	RespondJSON(w, 200, &Response{Message: "Updated"})
}
