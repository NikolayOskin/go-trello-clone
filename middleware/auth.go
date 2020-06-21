package middleware

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func Verified(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var u model.User
		jwtUser := r.Context().Value(UserCtx).(model.User)
		col := mongodb.Client.Database("trello").Collection("users")
		if err := col.FindOne(context.TODO(), bson.M{"_id": jwtUser.ID}).Decode(&u); err != nil {
			render.JSON(w, r, render.M{"error": err.Error()})
			return
		}
		if u.Verified == false {
			render.JSON(w, r, render.M{"error": "user must be verified to access this area"})
			return
		}

		next.ServeHTTP(w, r)
	})
}