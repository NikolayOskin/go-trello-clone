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

type CardController struct{}

func (c *CardController) Create(w http.ResponseWriter, r *http.Request) {
	var card model.Card
	user := r.Context().Value(mid.UserCtx).(model.User)
	if err := decodeJSON(w, r, &card); err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Could not parse json request"})
		return
	}
	card.UserId = user.ID.Hex()
	cardId, err := handlers.CreateCard(card)
	if err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Server error"})
		return
	}
	JSONResp(w, 200, &CreatedResponse{Message: "Added", Id: cardId})
}

func (c *CardController) Update(w http.ResponseWriter, r *http.Request) {
	card := r.Context().Value(mid.CardCtx).(model.Card)
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Server error"})
		return
	}
	card.ID = id
	if err = handlers.UpdateCard(card); err != nil {
		JSONResp(w, 200, &ErrResp{Message: "Server error"})
		return
	}
	JSONResp(w, 200, &Response{Message: "Updated"})
}

func (c *CardController) Delete(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(mid.UserCtx).(model.User)
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Server error"})
		return
	}
	col := mongodb.Client.Database("trello").Collection("cards")
	filter := bson.M{"_id": id, "user_id": user.ID.Hex()}
	if _, err = col.DeleteOne(context.TODO(), filter); err != nil {
		JSONResp(w, 500, &ErrResp{Message: err.Error()})
		return
	}
	JSONResp(w, 200, &Response{Message: "Deleted"})
}
