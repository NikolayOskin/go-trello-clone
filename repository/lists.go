package repository

import (
	"context"
	"time"

	"github.com/NikolayOskin/go-trello-clone/db"
	"github.com/NikolayOskin/go-trello-clone/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type listRepository struct{}

var Lists listRepository

func (l listRepository) GetById(ctx context.Context, id string) (*model.List, error) {
	var list model.List
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Lists.FindOne(ctx, bson.M{"_id": objId}).Decode(&list); err != nil {
		return nil, err
	}

	return &list, nil
}

func (l listRepository) GetByBoardId(ctx context.Context, id string) ([]model.List, error) {
	var lists []model.List
	filter := bson.D{{"board_id", id}}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := db.Lists.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &lists); err != nil {
		return nil, err
	}

	return lists, nil
}
