package handlers

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/NikolayOskin/go-trello-clone/db"
	pb "github.com/NikolayOskin/go-trello-clone/mailer/src"
	"github.com/NikolayOskin/go-trello-clone/model"
	mailer "github.com/NikolayOskin/go-trello-clone/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(reqUser *model.User, ctx context.Context) error {
	var user model.User
	filter := bson.M{"email": reqUser.Email}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Users.FindOne(ctx, filter).Decode(&user); err != nil {
		return errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password)); err != nil {
		return errors.New("invalid credentials")
	}
	reqUser.ID = user.ID

	return nil
}

func VerifyEmail(u model.User, code string, ctx context.Context) error {
	var user model.User

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Users.FindOne(ctx, bson.M{"_id": u.ID}).Decode(&user); err != nil {
		return err
	}
	if user.VerificationCode != code {
		return errors.New("incorrect verification code")
	}
	user.Verify()

	if _, err := db.Users.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{"$set": user}); err != nil {
		return err
	}

	return nil
}

func ResetPassword(email string, ctx context.Context) error {
	var u model.User

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Users.FindOne(ctx, bson.M{"email": email}).Decode(&u); err != nil {
		return errors.New("user with this email does not exist")
	}
	u.ResetPasswordCode = primitive.NewObjectID().Hex()
	_, err := db.Users.UpdateOne(ctx, bson.M{"email": email},
		bson.D{
			{"$set", bson.D{
				{"reset_password_code", u.ResetPasswordCode},
				{"reset_password_expired_at", time.Now().Add(time.Minute * 15)}},
			},
		})
	if err != nil {
		return err
	}

	if os.Getenv("APP_ENV") != "test" {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		_, err = mailer.Client().ResetPasswordEmail(ctx, &pb.EmailRequest{Email: email, Code: u.ResetPasswordCode})
		if err != nil {
			return err
		}
	}

	return nil
}

func SetNewPassword(email string, code string, password string) error {
	var u model.User
	if err := db.Users.FindOne(context.TODO(), bson.M{"email": email}).Decode(&u); err != nil {
		return errors.New("user with this email does not exist")
	}
	if u.ResetPasswordCode != code {
		return errors.New("invalid credentials")
	}
	if time.Now().After(u.ResetPasswordCodeExpiredAt) {
		return errors.New("code is expired")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	_, err = db.Users.UpdateOne(context.TODO(), bson.M{"email": email},
		bson.D{
			{"$set", bson.D{{"password", string(hash)}}},
		})
	if err != nil {
		return err
	}

	return nil
}
