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
	"go.mongodb.org/mongo-driver/bson"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func CreateUser(user model.User) error {
	var u model.User
	validate := v.New()
	if err := validate.Struct(user); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return errors.New(e.Translate(v.Trans))
		}
	}
	if err := mongodb.Users.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&u); err == nil {
		return errors.New("user with this email already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)
	user.Verified = false
	user.VerificationCode = rand.Int()

	if _, err = mongodb.Users.InsertOne(context.TODO(), user); err != nil {
		return err
	}

	if err = SendCode(user.VerificationCode, user.Email); err != nil {
		fmt.Println(err)
	}

	return nil
}

func SendCode(code int, email string) error {
	conn, err := grpc.Dial(os.Getenv("MAILER_SERVICE_PORT"), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	c := pb.NewMailerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err = c.SendEmail(ctx, &pb.EmailRequest{Email: email, Code: strconv.Itoa(code)})
	if err != nil {
		return err
	}

	return nil
}
