package controller

import (
	"github.com/NikolayOskin/go-trello-clone/handlers"
	mid "github.com/NikolayOskin/go-trello-clone/middleware"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/repository"

	//"github.com/NikolayOskin/go-trello-clone/repository"

	//"github.com/NikolayOskin/go-trello-clone/repository"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type BoardController struct{}

func (b *BoardController) GetFull(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(mid.UserCtx).(model.User)
	boardRepo := repository.Boards{}
	board, err := boardRepo.GetById(chi.URLParam(r, "id"))
	if err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Server error"})
		return
	}
	if board == nil || board.UserId != userCtx.ID.Hex() {
		JSONResp(w, 404, &ErrResp{Message: "Not found"})
		return
	}
	cardsRepo := repository.Cards{}
	cards, err := cardsRepo.GetByBoardId(board.ID.Hex())
	if err != nil {
		JSONResp(w, 500, &ErrResp{Message: err.Error()})
		return
	}

	m := make(map[string][]model.Card)
	for _, card := range cards {
		if _, ok := m[card.ListId]; ok == true {
			m[card.ListId] = append(m[card.ListId], card)
		} else {
			m[card.ListId] = []model.Card{card}
		}
	}
	listsRepo := repository.Lists{}
	lists, err := listsRepo.GetByBoardId(board.ID.Hex())
	if err != nil {
		JSONResp(w, 500, &ErrResp{Message: err.Error()})
		return
	}
	for i, list := range lists {
		lists[i].Cards = m[list.ID.Hex()]
	}
	board.Lists = lists

	JSONResp(w, 200, board)
}

func (b *BoardController) Create(w http.ResponseWriter, r *http.Request) {
	var board model.Board
	user := r.Context().Value(mid.UserCtx).(model.User)
	if err := decodeJSON(w, r, &board); err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Could not parse json request"})
		return
	}
	board.UserId = user.ID.Hex()
	if err := handlers.CreateBoard(board); err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Server error"})
		return
	}
	JSONResp(w, 200, &Response{Message: "Added"})
}

func (b *BoardController) Update(w http.ResponseWriter, r *http.Request) {
	board := r.Context().Value(mid.BoardCtx).(model.Board)
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Server error"})
		return
	}
	board.ID = id
	if err := handlers.UpdateBoard(board); err != nil {
		JSONResp(w, 200, &ErrResp{Message: "Server error"})
		return
	}
	JSONResp(w, 200, &Response{Message: "Updated"})
}
