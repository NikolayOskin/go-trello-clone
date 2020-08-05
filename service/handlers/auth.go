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
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(ctx context.Context, reqUser *model.User) error {
	var user model.User
	filter := bson.M{"email": reqUser.Email}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Users.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("invalid credentials")
		}
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password)); err != nil {
		return errors.New("invalid credentials")
	}
	reqUser.ID = user.ID

	return nil
}

func VerifyEmail(ctx context.Context, u model.User, code string) error {
	var user model.User

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Users.FindOne(ctx, bson.M{"_id": u.ID}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user not found")
		}
		return err
	}
	if user.VerificationCode != code {
		return errors.New("incorrect verification code")
	}
	user.Verify()

	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := db.Users.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{"$set": user}); err != nil {
		return err
	}

	return nil
}

func ResetPassword(ctx context.Context, email string) error {
	var u model.User

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Users.FindOne(ctx, bson.M{"email": email}).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user not found")
		}
		return err
	}
	u.ResetPasswordCode = primitive.NewObjectID().Hex()

	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

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

func SetNewPassword(ctx context.Context, email string, code string, password string) error {
	var u model.User

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Users.FindOne(ctx, bson.M{"email": email}).Decode(&u); err != nil {
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

	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = db.Users.UpdateOne(ctx, bson.M{"email": email},
		bson.D{
			{"$set", bson.D{{"password", string(hash)}}},
		})
	if err != nil {
		return err
	}

	return nil
}
