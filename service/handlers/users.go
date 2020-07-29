package handlers

import (
	"context"
	"errors"
	"github.com/NikolayOskin/go-trello-clone/db"
	pb "github.com/NikolayOskin/go-trello-clone/mailer/src"
	"github.com/NikolayOskin/go-trello-clone/model"
	mailer "github.com/NikolayOskin/go-trello-clone/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"

	"golang.org/x/crypto/bcrypt"
	"time"
)

func CreateUser(ctx context.Context, user model.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	user.Verified = false
	user.VerificationCode = primitive.NewObjectID().Hex()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err = db.Users.InsertOne(ctx, user); err != nil {
		if db.IsDuplicated(err) {
			return errors.New("user with this email already exists")
		}
		return err
	}

	if os.Getenv("APP_ENV") != "test" {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		_, err = mailer.Client().SignUpEmail(ctx, &pb.EmailRequest{Email: user.Email, Code: user.VerificationCode})
		if err != nil {
			return err
		}
	}

	return nil
}
