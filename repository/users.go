package repository

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct{}

func (b *Users) GetById(id string) (*model.User, error) {
	var user model.User
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := mongodb.Users.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
