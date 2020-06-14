package handlers

import (
	"context"
	"fmt"
	pb "github.com/NikolayOskin/go-trello-clone/mailer/mailerpkg"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func CreateUser(user model.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)
	user.Verified = false
	user.VerificationCode = rand.Int()

	col := mongodb.Client.Database("trello").Collection("users")
	_, err = col.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	err = sendCode(user.VerificationCode, user.Email)
	if err != nil {
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
