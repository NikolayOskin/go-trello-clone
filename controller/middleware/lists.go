package middleware

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	v "github.com/NikolayOskin/go-trello-clone/service/validator"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func DecodeListObj(validate *validator.Validate) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			list := model.UpdateList{}
			currUser := r.Context().Value(UserCtx).(model.User)

			if err := render.DecodeJSON(r.Body, &list); err != nil {
				w.WriteHeader(400)
				render.JSON(w, r, render.M{"error": err.Error()})
				return
			}

			if err := validate.Struct(list); err != nil {
				for _, e := range err.(validator.ValidationErrors) {
					w.WriteHeader(422)
					render.JSON(w, r, render.M{"error": e.Translate(v.Trans)})
					return
				}
			}

			list.UserId = currUser.ID.Hex()
			ctx := context.WithValue(r.Context(), ListCtx, list)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
