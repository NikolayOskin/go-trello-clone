package middleware

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/go-chi/render"
	"net/http"
)

func DecodeBoardObj(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		board := model.Board{}
		currentUser := r.Context().Value(UserCtx).(model.User)

		err := render.DecodeJSON(r.Body, &board)
		if err != nil {
			render.JSON(w, r, render.M{"error": err.Error()})
			return
		}

		board.UserId = currentUser.ID.Hex()

		ctx := context.WithValue(r.Context(), BoardCtx, board)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
