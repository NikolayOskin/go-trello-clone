package handlers

import (
	"context"
	"fmt"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"golang.org/x/crypto/bcrypt"
)

func HandleAddUser(user model.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)

	collection := mongodb.Client.Database("trello").Collection("users")

	insertedUser, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	fmt.Println("inserted user: ", insertedUser.InsertedID)
	return nil
}