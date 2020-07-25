package repository

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/db"
	"github.com/NikolayOskin/go-trello-clone/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cards struct{}

func (c *Cards) GetById(id string) (*model.Card, error) {
	var card model.Card
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := db.Cards.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&card); err != nil {
		return nil, err
	}

	return &card, nil
}

func (c *Cards) GetByBoardId(id string) ([]model.Card, error) {
	var cards []model.Card
	filter := bson.D{{"board_id", id}}
	cursor, err := db.Cards.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &cards); err != nil {
		return nil, err
	}

	return cards, nil
}
