package middleware

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	v "github.com/NikolayOskin/go-trello-clone/service/validator"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func DecodeBoardObj(validate *validator.Validate) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			board := model.Board{}
			currUser := r.Context().Value(UserCtx).(model.User)

			if err := render.DecodeJSON(r.Body, &board); err != nil {
				w.WriteHeader(400)
				render.JSON(w, r, render.M{"error": err.Error()})
				return
			}

			if err := validate.Struct(board); err != nil {
				for _, e := range err.(validator.ValidationErrors) {
					w.WriteHeader(422)
					render.JSON(w, r, render.M{"error": e.Translate(v.Trans)})
					return
				}
			}

			board.UserId = currUser.ID.Hex()
			ctx := context.WithValue(r.Context(), BoardCtx, board)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
