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
	col := mongodb.Client.Database("trello").Collection("cards")
	filter := bson.D{{"board_id", id}}
	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &cards); err != nil {
		return nil, err
	}

	return cards, nil
}
