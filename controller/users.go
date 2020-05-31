package controller

import (
	"github.com/NikolayOskin/go-trello-clone/handlers"
	"github.com/NikolayOskin/go-trello-clone/model"
	"net/http"
)

func AddUser(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := decodeJSONBody(w, r, &user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = handlers.HandleAddUser(user)
	if err != nil {
		RespondJSON(w, 400, &Response{Message:err.Error()})
		return
	}

	RespondJSON(w, 200, &Response{Message:"Success"})
}

func GetAuthUser(w http.ResponseWriter, r *http.Request) {

}