package controller

import (
	"net/http"

	mid "github.com/NikolayOskin/go-trello-clone/controller/middleware"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/repository"
	"github.com/NikolayOskin/go-trello-clone/service/handlers"
	v "github.com/NikolayOskin/go-trello-clone/service/validator"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardController struct {
	Validate     *validator.Validate
	BoardHandler handlers.Board
}

func (b *BoardController) GetFull(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(mid.UserCtx).(model.User)
	board, err := repository.Boards.GetById(r.Context(), chi.URLParam(r, "id"))
	if err != nil || board == nil || board.UserId != userCtx.ID.Hex() {
		JSONResp(w, 404, ErrResp{Message: "Not found"})
		return
	}
	if err := b.BoardHandler.FillBoardWithListsAndCards(r.Context(), board); err != nil {
		JSONResp(w, 500, ErrResp{Message: err.Error()})
		return
	}
	JSONResp(w, 200, board)
}

func (b *BoardController) Create(w http.ResponseWriter, r *http.Request) {
	var board model.Board
	user := r.Context().Value(mid.UserCtx).(model.User)
	if err := decodeJSON(w, r, &board); err != nil {
		JSONResp(w, err.(*malformedRequest).Status, ErrResp{err.Error()})
		return
	}
	if err := b.Validate.Struct(board); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			JSONResp(w, 422, ErrResp{e.Translate(v.Trans)})
			return
		}
	}
	board.UserId = user.ID.Hex()
	boardId, err := b.BoardHandler.Create(r.Context(), board)
	if err != nil {
		JSONResp(w, 500, ErrResp{Message: "Server error"})
		return
	}
	JSONResp(w, 201, CreatedResponse{Message: "Created", Id: boardId})
}

func (b *BoardController) Update(w http.ResponseWriter, r *http.Request) {
	board := r.Context().Value(mid.BoardCtx).(model.Board)
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		JSONResp(w, 500, ErrResp{Message: "Server error"})
		return
	}
	board.ID = id
	if err := b.BoardHandler.Update(r.Context(), board); err != nil {
		JSONResp(w, 500, ErrResp{Message: "Server error"})
		return
	}
	JSONResp(w, 200, Response{Message: "Updated"})
}
