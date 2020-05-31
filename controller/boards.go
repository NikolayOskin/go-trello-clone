package controller

import (
	"github.com/NikolayOskin/go-trello-clone/handlers"
	"github.com/NikolayOskin/go-trello-clone/model"
	"net/http"
)

func AddBoard(w http.ResponseWriter, r *http.Request) {
	var board model.Board

	err := decodeJSONBody(w, r, &board)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	handlers.HandleAddBoard(board)

	RespondJSON(w, 200, &Response{Message:"Success"})
}

func GetBoards(w http.ResponseWriter, r *http.Request) {

}