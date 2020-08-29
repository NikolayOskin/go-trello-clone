package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/NikolayOskin/go-trello-clone/db"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Card struct{}

func (h Card) Create(ctx context.Context, c model.Card) (string, error) {
	list, err := repository.Lists.GetById(ctx, c.ListId)
	if err != nil {
		return "", err
	}
	if list == nil || list.UserId != c.UserId {
		return "", errors.New("list for your user does not exist")
	}
	c.BoardId = list.BoardId
	res, err := db.Cards.InsertOne(context.TODO(), c)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (h Card) Update(ctx context.Context, c model.Card) error {
	card, err := repository.Cards.GetById(ctx, c.ID.Hex())
	if err != nil {
		return err
	}

	list, err := repository.Lists.GetById(ctx, c.ListId)
	if err != nil {
		return err
	}

	// check if user tries to update card which belongs to another user
	if card.UserId != c.UserId || c.BoardId != card.BoardId || list.UserId != c.UserId {
		return errors.New("board or list does not exist")
	}

	filter := bson.M{"_id": c.ID}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := db.Cards.UpdateOne(ctx, filter, bson.M{"$set": c}); err != nil {
		return err
	}
	return nil
}
