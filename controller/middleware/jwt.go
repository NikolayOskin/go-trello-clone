package middleware

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/service/auth"
	"net/http"

	"github.com/go-chi/render"
)

func JWTCheck(auth *auth.Auth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var ctx context.Context
			tokenString := r.Header.Get("Authorization")

			if len(tokenString) == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, render.M{"message": "Unauthenticated"})
				return
			}

			claims, err := auth.ValidateToken(tokenString)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, render.M{"message": "Unauthenticated"})
				return
			}

			ctx = context.WithValue(r.Context(), UserCtx, claims.User)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
