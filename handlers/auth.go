package handlers

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func HandleAuthenticate(requested *model.User) error {
	var user model.User
	col := mongodb.Client.Database("trello").Collection("users")
	filter := bson.M{"email": requested.Email}
	if err := col.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requested.Password)); err != nil {
		return err
	}
	requested.ID = user.ID

	return nil
}
