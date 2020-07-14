package mongodb

import (
	"context"
	"errors"
	"github.com/NikolayOskin/go-trello-clone/model"
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
	seedUser(verifiedUser)
	seedUser(userWithPassResetCode)
	seedUser(userWithExpiredResetCode)
}

func seedUser(user *model.User) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = string(hash)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err = Users.InsertOne(ctx, user); err != nil {
		if IsDuplicated(err) {
			log.Fatal(errors.New("user with this email already exists"))
		}
	}
}
