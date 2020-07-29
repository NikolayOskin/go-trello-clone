package controller

import (
	"context"
	mid "github.com/NikolayOskin/go-trello-clone/controller/middleware"
	"github.com/NikolayOskin/go-trello-clone/db"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/service/handlers"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type CardController struct{}

func (c *CardController) Create(w http.ResponseWriter, r *http.Request) {
	var card model.Card
	user := r.Context().Value(mid.UserCtx).(model.User)
	if err := decodeJSON(w, r, &card); err != nil {
		JSONResp(w, err.(*malformedRequest).Status, &ErrResp{err.Error()})
		return
	}
	card.UserId = user.ID.Hex()
	cardId, err := handlers.CreateCard(r.Context(), card)
	if err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Server error"})
		return
	}
	JSONResp(w, 201, &CreatedResponse{Message: "Created", Id: cardId})
}

func (c *CardController) Update(w http.ResponseWriter, r *http.Request) {
	card := r.Context().Value(mid.CardCtx).(model.Card)
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Server error"})
		return
	}
	card.ID = id
	if err = handlers.UpdateCard(r.Context(), card); err != nil {
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
	filter := bson.M{"_id": id, "user_id": user.ID.Hex()}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if _, err = db.Cards.DeleteOne(ctx, filter); err != nil {
		JSONResp(w, 500, &ErrResp{Message: err.Error()})
		return
	}
	JSONResp(w, 200, &Response{Message: "Deleted"})
}
