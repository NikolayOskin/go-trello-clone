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

	err := decodeJSONBody(w, r, &card)
	if err != nil {
		RespondJSON(w, 500, &ErrResp{Message: "Could not parse json request"})
		return
	}
	card.UserId = user.ID.Hex()
	err = handlers.HandleCreateCard(card)
	if err != nil {
		RespondJSON(w, 500, &ErrResp{Message: "Server error"})
		return
	}

	RespondJSON(w, 200, &Response{Message: "Added"})
}

func (c *CardController) Delete(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(mid.UserCtx).(model.User)
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		RespondJSON(w, 500, &ErrResp{Message: "Server error"})
		return
	}
	col := mongodb.Client.Database("trello").Collection("cards")
	_, err = col.DeleteOne(context.TODO(), bson.M{"_id": id, "user_id": user.ID.Hex()})
	if err != nil {
		RespondJSON(w, 500, &ErrResp{Message: err.Error()})
		return
	}

	RespondJSON(w, 200, &Response{Message: "Deleted"})
}

func (c *CardController) Update(w http.ResponseWriter, r *http.Request) {
	card := r.Context().Value(mid.CardCtx).(model.Card)

	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		RespondJSON(w, 500, &ErrResp{Message: "Server error"})
		return
	}
	card.ID = id
	err = handlers.HandleUpdateCard(card)
	if err != nil {
		RespondJSON(w, 200, &ErrResp{Message: "Server error"})
		return
	}

	RespondJSON(w, 200, &Response{Message: "Updated"})
}
