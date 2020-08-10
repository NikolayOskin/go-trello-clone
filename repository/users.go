package repository

import (
	"context"
	"errors"
	"github.com/NikolayOskin/go-trello-clone/db"
	"github.com/NikolayOskin/go-trello-clone/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Users struct{}

func (b *Users) FindById(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Users.FindOne(ctx, bson.M{"_id": objId}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (b *Users) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Users.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
