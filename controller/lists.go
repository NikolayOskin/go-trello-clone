package controller

import (
	mid "github.com/NikolayOskin/go-trello-clone/controller/middleware"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/service/handlers"
	v "github.com/NikolayOskin/go-trello-clone/service/validator"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type ListController struct {
	Validate *validator.Validate
}

func (l *ListController) Create(w http.ResponseWriter, r *http.Request) {
	var list model.List
	user := r.Context().Value(mid.UserCtx).(model.User)
	if err := decodeJSON(w, r, &list); err != nil {
		JSONResp(w, err.(*malformedRequest).Status, &ErrResp{err.Error()})
		return
	}
	list.UserId = user.ID.Hex()
	if err := l.Validate.Struct(list); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			JSONResp(w, 422, &ErrResp{e.Translate(v.Trans)})
			return
		}
	}
	listId, err := handlers.CreateList(list)
	if err != nil {
		JSONResp(w, 400, &ErrResp{Message: err.Error()})
		return
	}
	JSONResp(w, 201, &CreatedResponse{Message: "Created", Id: listId})
}

func (l *ListController) Update(w http.ResponseWriter, r *http.Request) {
	list := r.Context().Value(mid.ListCtx).(model.List)
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		JSONResp(w, 500, &ErrResp{Message: "Server error"})
		return
	}
	list.ID = id
	if err = handlers.UpdateList(list); err != nil {
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
	if err := handlers.DeleteList(id, user); err != nil {
		JSONResp(w, 500, &ErrResp{Message: err.Error()})
		return
	}
	JSONResp(w, 200, &Response{Message: "Deleted"})
}
