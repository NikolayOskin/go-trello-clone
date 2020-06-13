package controller

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/handlers"
	mid "github.com/NikolayOskin/go-trello-clone/middleware"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type ListController struct{}

func (l *ListController) Create(w http.ResponseWriter, r *http.Request) {
	var list model.List
	user := r.Context().Value(mid.UserCtx).(model.User)
	if err := decodeJSON(w, r, &list); err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Could not parse json request"})
		return
	}
	list.UserId = user.ID.Hex()
	if err := handlers.HandleCreateList(list); err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Server error"})
		return
	}
	JSONResp(w, 200, &Response{Message: "Added"})
}

func (l *ListController) Update(w http.ResponseWriter, r *http.Request) {
	list := r.Context().Value(mid.ListCtx).(model.List)
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Server error"})
		return
	}
	list.ID = id
	if err = handlers.HandleUpdateList(list); err != nil {
		JSONResp(w, 200, &ErrResp{Message: "Server error"})
		return
	}
	JSONResp(w, 200, &Response{Message: "Updated"})
}

func (l *ListController) Delete(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(mid.UserCtx).(model.User)
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Server error"})
		return
	}
	// first deleting cards associated with list
	cardsCol := mongodb.Client.Database("trello").Collection("cards")
	f := bson.M{"list_id": id.Hex(), "user_id": user.ID.Hex()}
	if _, err := cardsCol.DeleteMany(context.TODO(), f); err != nil {
		JSONResp(w, 500, &ErrResp{Message: err.Error()})
		return
	}
	// then delete list itself
	listsCol := mongodb.Client.Database("trello").Collection("lists")
	f = bson.M{"_id": id, "user_id": user.ID.Hex()}
	if _, err := listsCol.DeleteOne(context.TODO(), f); err != nil {
		JSONResp(w, 500, &ErrResp{Message: err.Error()})
		return
	}
	JSONResp(w, 200, &Response{Message: "Deleted"})
}
