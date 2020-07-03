package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	v "github.com/NikolayOskin/go-trello-clone/service/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(reqUser *model.User) error {
	validate := v.New()
	if err := validate.Struct(reqUser); err != nil {
		return err
	}
	var user model.User
	filter := bson.M{"email": reqUser.Email}
	if err := mongodb.Users.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		return errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password)); err != nil {
		return errors.New("invalid credentials")
	}
	reqUser.ID = user.ID

	return nil
}

func VerifyEmail(u model.User, code string) error {
	var user model.User
	if code == "" {
		return errors.New("code must not be empty")
	}
	if err := mongodb.Users.FindOne(context.TODO(), bson.M{"_id": u.ID}).Decode(&user); err != nil {
		return err
	}
	if user.VerificationCode != code {
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
	u.ResetPasswordCode = primitive.NewObjectID().Hex()
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

func SetNewPassword(email string, code string, password string) error {
	var u model.User
	if err := mongodb.Users.FindOne(context.TODO(), bson.M{"email": email}).Decode(&u); err != nil {
		return errors.New("user with this email does not exist")
	}
	if u.ResetPasswordCode != code {
		return errors.New("invalid credentials")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	_, err = mongodb.Users.UpdateOne(context.TODO(), bson.M{"email": email},
		bson.D{
			{"$set", bson.D{{"password", string(hash)}}},
		})
	if err != nil {
		return err
	}
	return nil
}
