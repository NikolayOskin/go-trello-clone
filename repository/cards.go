package repository

import (
	"context"
	"time"

	"github.com/NikolayOskin/go-trello-clone/db"
	"github.com/NikolayOskin/go-trello-clone/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type cardRepository struct{}

var Cards cardRepository

func (c cardRepository) GetById(ctx context.Context, id string) (*model.Card, error) {
	var card model.Card
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Cards.FindOne(ctx, bson.M{"_id": objId}).Decode(&card); err != nil {
		return nil, err
	}

	return &card, nil
}

func (c cardRepository) GetByBoardId(ctx context.Context, id string) ([]model.Card, error) {
	var cards []model.Card
	filter := bson.D{{"board_id", id}}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := db.Cards.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &cards); err != nil {
		return nil, err
	}

	return cards, nil
}
