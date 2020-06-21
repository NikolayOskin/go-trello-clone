package handlers

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/NikolayOskin/go-trello-clone/mailer/mailerpkg"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	v "github.com/NikolayOskin/go-trello-clone/validator"
	"github.com/go-playground/validator/v10"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func CreateUser(user model.User) error {
	validate := v.New()
	if err := validate.Struct(user); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return errors.New(e.Translate(v.Trans))
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)
	user.Verified = false
	user.VerificationCode = rand.Int()

	col := mongodb.Client.Database("trello").Collection("users")
	if _, err = col.InsertOne(context.TODO(), user); err != nil {
		return err
	}

	if err = sendCode(user.VerificationCode, user.Email); err != nil {
		fmt.Println(err)
	}

	return nil
}

func sendCode(code int, email string) error {
	conn, err := grpc.Dial(os.Getenv("MAILER_SERVICE_PORT"), grpc.WithInsecure())
	if err != nil {
		return err
	}
	c := pb.NewMailerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err = c.SendEmail(ctx, &pb.EmailRequest{Email: email, Code: strconv.Itoa(code)})
	if err != nil {
		return err
	}

	return nil
}
