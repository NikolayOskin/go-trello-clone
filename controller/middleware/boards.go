package middleware

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/service/validator"
	"github.com/go-chi/render"
	"net/http"
)

func DecodeBoardObj(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := model.Board{}
		currUser := r.Context().Value(UserCtx).(model.User)
		if err := render.DecodeJSON(r.Body, &b); err != nil {
			render.JSON(w, r, render.M{"error": err.Error()})
			return
		}
		validate := validator.New()
		if err := validate.Struct(b); err != nil {
			render.JSON(w, r, render.M{"error": err.Error()})
			return
		}
		b.UserId = currUser.ID.Hex()
		ctx := context.WithValue(r.Context(), BoardCtx, b)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
