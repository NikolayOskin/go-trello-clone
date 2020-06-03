package middleware

import (
	"context"
	"errors"
	"github.com/NikolayOskin/go-trello-clone/config"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

func JWTCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context
		tokenString := r.Header.Get("Authorization")

		if len(tokenString) == 0 {
			http.Error(w, "Unauthenticated", 401)
			return
		}

		cl := model.JWTClaims{}

		_, err := jwt.ParseWithClaims(tokenString, &cl, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(config.JWTSecret), nil
		})

		if err != nil {
			http.Error(w, "Unauthenticated", 401)
			return
		}
		ctx = context.WithValue(r.Context(), UserCtx, cl.User)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
