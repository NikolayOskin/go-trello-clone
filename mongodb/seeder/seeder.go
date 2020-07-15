package seeder

import (
	"context"
	"errors"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

func Seed() {
	if os.Getenv("APP_ENV") != "test" {
		log.Fatal("you can seed db only in test environment")
	}
	unverifiedUser := &model.User{
		Email:            "testuser@gmail.com",
		Password:         "qwerty",
		Verified:         false,
		VerificationCode: "12345",
	}
	verifiedUser := &model.User{
		Email:    "verified@gmail.com",
		Password: "qwerty",
		Verified: true,
	}
	userWithPassResetCode := &model.User{
		Email:                      "withresetpasswordcode@gmail.com",
		Password:                   "qwerty",
		Verified:                   true,
		ResetPasswordCode:          "123456",
		ResetPasswordCodeExpiredAt: time.Now().Add(15 * time.Minute),
	}
	userWithExpiredResetCode := &model.User{
		Email:                      "withexpiredresetpasswordcode@gmail.com",
		Password:                   "qwerty",
		Verified:                   true,
		ResetPasswordCode:          "123456",
		ResetPasswordCodeExpiredAt: time.Now().AddDate(0, 0, -1),
	}
	seedUser(unverifiedUser)
	verifiedUserId := seedUser(verifiedUser)
	seedUser(userWithPassResetCode)
	seedUser(userWithExpiredResetCode)

	boardId := seedBoard("507f191e810c19729de860ea", "Some board", verifiedUserId)
	listId := seedList("507f191e810c19729de860ea", "First list", boardId, verifiedUserId)
	seedCard("507f191e810c19729de860ea", "first card", listId, boardId, verifiedUserId)
}

func seedUser(user *model.User) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = string(hash)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := mongodb.Users.InsertOne(ctx, user)
	if err != nil {
		if mongodb.IsDuplicated(err) {
			log.Fatal(errors.New("user with this email already exists"))
		}
		log.Fatal(err)
	}

	return res.InsertedID.(primitive.ObjectID).Hex()
}

func seedBoard(boardHex string, title string, userId string) string {
	objId, _ := primitive.ObjectIDFromHex(boardHex)
	b := model.Board{
		ID:     objId,
		Title:  title,
		UserId: userId,
		Lists:  nil,
	}
	res, err := mongodb.Boards.InsertOne(context.Background(), b)
	if err != nil {
		log.Fatal(err)
	}
	return res.InsertedID.(primitive.ObjectID).Hex()
}

func seedList(listHex string, title string, boardId string, userId string) string {
	objId, _ := primitive.ObjectIDFromHex(listHex)
	l := model.List{
		ID:       objId,
		Title:    title,
		BoardId:  boardId,
		UserId:   userId,
		Position: 1,
	}
	res, err := mongodb.Lists.InsertOne(context.Background(), l)
	if err != nil {
		log.Fatal(err)
	}
	return res.InsertedID.(primitive.ObjectID).Hex()
}

func seedCard(cardHex string, text string, listId string, boardId string, userId string) string {
	objId, _ := primitive.ObjectIDFromHex(cardHex)
	c := model.Card{
		ID:       objId,
		Text:     text,
		BoardId:  boardId,
		ListId:   listId,
		UserId:   userId,
		Position: 1,
	}
	res, err := mongodb.Cards.InsertOne(context.Background(), c)
	if err != nil {
		log.Fatal(err)
	}
	return res.InsertedID.(primitive.ObjectID).Hex()
}
