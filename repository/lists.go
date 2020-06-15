package repository

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lists struct{}

func (l *Lists) GetById(id string) (*model.List, error) {
	var list model.List
	col := mongodb.Client.Database("trello").Collection("lists")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := col.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&list); err != nil {
		return nil, err
	}

	return &list, nil
}
