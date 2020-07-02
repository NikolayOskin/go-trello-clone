package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	v "github.com/NikolayOskin/go-trello-clone/validator"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
)

func Authenticate(requested *model.User) error {
	validate := v.New()
	if err := validate.Struct(requested); err != nil {
		return err
	}
	var user model.User
	filter := bson.M{"email": requested.Email}
	if err := mongodb.Users.FindOne(context.TODO(), filter).Decode(&user); err != nil {
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
	if err := mongodb.Users.FindOne(context.TODO(), bson.M{"_id": u.ID}).Decode(&user); err != nil {
		return err
	}
	if user.VerificationCode != rCode {
		return errors.New("incorrect verification code")
	}
	user.Verified = true
	if _, err := mongodb.Users.UpdateOne(context.TODO(), bson.M{"_id": u.ID}, bson.M{"$set": user}); err != nil {
		return err
	}
	return nil
}

func ResetPassword(email string) error {
	var u model.User
	if err := mongodb.Users.FindOne(context.TODO(), bson.M{"email": email}).Decode(&u); err != nil {
		return errors.New("user with this email does not exist")
	}
	u.ResetPasswordCode = rand.Int()
	_, err := mongodb.Users.UpdateOne(context.TODO(), bson.M{"email": email},
		bson.D{
			{"$set", bson.D{{"reset_password_code", u.ResetPasswordCode}}},
		})
	if err != nil {
		return err
	}
	if err = SendCode(u.ResetPasswordCode, email); err != nil {
		fmt.Println(err)
	}

	return nil
}
