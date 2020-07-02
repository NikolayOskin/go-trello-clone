package repository

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

type Cards struct{}

func (c *Cards) GetByBoardId(id string) ([]model.Card, error) {
	var cards []model.Card
	filter := bson.D{{"board_id", id}}
	cursor, err := mongodb.Cards.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &cards); err != nil {
		return nil, err
	}

	return cards, nil
}
