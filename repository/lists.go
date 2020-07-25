package repository

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/db"
	"github.com/NikolayOskin/go-trello-clone/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lists struct{}

func (l *Lists) GetById(id string) (*model.List, error) {
	var list model.List
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := db.Lists.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&list); err != nil {
		return nil, err
	}

	return &list, nil
}

func (l *Lists) GetByBoardId(id string) ([]model.List, error) {
	var lists []model.List
	filter := bson.D{{"board_id", id}}
	cursor, err := db.Lists.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &lists); err != nil {
		return nil, err
	}

	return lists, nil
}
