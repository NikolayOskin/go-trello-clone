package handlers

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(requested *model.User) error {
	validate := validator.New()
	if err := validate.Struct(requested); err != nil {
		return err
	}
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
