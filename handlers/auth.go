package handlers

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func HandleAuthenticate(userRequest *model.User) error {
	var user model.User

	collection := mongodb.Client.Database("trello").Collection("users")
	filter := bson.M{"email": userRequest.Email}

	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		return err
	}

	userRequest.ID = user.ID

	return nil
}
