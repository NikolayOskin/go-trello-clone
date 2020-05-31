package handlers

import (
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func HandleAuthenticate(userRequest model.User) error {
	var user model.User

	collection := mongodb.Client.Database("trello").Collection("users")
	filter := bson.D{{"email", userRequest.Email}}

	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		return err
	}

	return nil
}
