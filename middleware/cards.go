package middleware

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func DecodeCardObj(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := model.Card{}
		currUser := r.Context().Value(UserCtx).(model.User)
		if err := render.DecodeJSON(r.Body, &c); err != nil {
			render.JSON(w, r, render.M{"error": err.Error()})
			return
		}
		validate := validator.New()
		if err := validate.Struct(c); err != nil {
			render.JSON(w, r, render.M{"error": err.Error()})
			return
		}
		c.UserId = currUser.ID.Hex()
		ctx := context.WithValue(r.Context(), CardCtx, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
