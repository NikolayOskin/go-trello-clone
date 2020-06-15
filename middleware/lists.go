package middleware

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func DecodeListObj(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := model.List{}
		currUser := r.Context().Value(UserCtx).(model.User)
		if err := render.DecodeJSON(r.Body, &l); err != nil {
			render.JSON(w, r, render.M{"error": err.Error()})
			return
		}
		validate := validator.New()
		if err := validate.Struct(l); err != nil {
			render.JSON(w, r, render.M{"error": err.Error()})
			return
		}
		l.UserId = currUser.ID.Hex()
		ctx := context.WithValue(r.Context(), ListCtx, l)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
