package handlers

import (
	"context"
	"errors"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	v "github.com/NikolayOskin/go-trello-clone/validator"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

func Authenticate(requested *model.User) error {
	validate := v.New()
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

func VerifyEmail(u model.User, code string) error {
	var user model.User
	if code == "" {
		return errors.New("code must not be empty")
	}
	rCode, err := strconv.Atoi(code)
	if err != nil {
		return err
	}
	col := mongodb.Client.Database("trello").Collection("users")
	if err := col.FindOne(context.TODO(), bson.M{"_id": u.ID}).Decode(&user); err != nil {
		return err
	}
	if user.VerificationCode != rCode {
		return errors.New("incorrect verification code")
	}
	user.Verified = true
	if _, err := col.UpdateOne(context.TODO(), bson.M{"_id": u.ID}, bson.M{"$set": user}); err != nil {
		return err
	}
	return nil
}
