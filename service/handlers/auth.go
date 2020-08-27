package handlers

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/NikolayOskin/go-trello-clone/db"
	pb "github.com/NikolayOskin/go-trello-clone/mailer/src"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/repository"
	mailer "github.com/NikolayOskin/go-trello-clone/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct{}

func (h Auth) Authenticate(ctx context.Context, reqUser *model.User) error {
	user, err := repository.Users.FindByEmail(ctx, reqUser.Email)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password)); err != nil {
		return errors.New("invalid credentials")
	}
	reqUser.ID = user.ID

	return nil
}

func (h Auth) VerifyEmail(ctx context.Context, u model.User, code string) error {
	user, err := repository.Users.FindById(ctx, u.ID.Hex())
	if err != nil {
		return err
	}
	if user.VerificationCode != code {
		return errors.New("incorrect verification code")
	}
	user.Verify()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := db.Users.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{"$set": *user}); err != nil {
		return err
	}

	return nil
}

func (h Auth) ResetPassword(ctx context.Context, email string) error {
	user, err := repository.Users.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	user.ResetPasswordCode = primitive.NewObjectID().Hex()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = db.Users.UpdateOne(ctx, bson.M{"email": email},
		bson.D{
			{"$set", bson.D{
				{"reset_password_code", user.ResetPasswordCode},
				{"reset_password_expired_at", time.Now().Add(time.Minute * 15)}},
			},
		})
	if err != nil {
		return err
	}

	if os.Getenv("APP_ENV") != "test" {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		_, err = mailer.Client().ResetPasswordEmail(ctx, &pb.EmailRequest{Email: email, Code: user.ResetPasswordCode})
		if err != nil {
			return err
		}
	}

	return nil
}

func (h Auth) SetNewPassword(ctx context.Context, email string, code string, password string) error {
	user, err := repository.Users.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user.ResetPasswordCode != code {
		return errors.New("invalid code")
	}
	if time.Now().After(user.ResetPasswordCodeExpiredAt) {
		return errors.New("code is expired")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = db.Users.UpdateOne(ctx, bson.M{"email": email},
		bson.D{
			{"$set", bson.D{{"password", string(hash)}}},
			{"$unset", bson.A{"reset_password_code", "reset_password_expired_at"}},
		})
	if err != nil {
		return err
	}

	return nil
}
