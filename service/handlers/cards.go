package handlers

import (
	"context"
	"errors"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"github.com/NikolayOskin/go-trello-clone/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCard(c model.Card) (string, error) {
	repo := repository.Lists{}
	list, err := repo.GetById(c.ListId)
	if err != nil {
		return "", err
	}
	if list == nil || list.UserId != c.UserId {
		return "", errors.New("list for your user does not exist")
	}
	c.BoardId = list.BoardId
	res, err := mongodb.Cards.InsertOne(context.TODO(), c)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func UpdateCard(c model.Card) error {
	cardsRepo := repository.Cards{}
	listsRepo := repository.Lists{}

	card, err := cardsRepo.GetById(c.ID.Hex())
	if err != nil {
		return err
	}

	list, err := listsRepo.GetById(c.ListId)
	if err != nil {
		return err
	}

	// check if user tries to update card which belongs to another user
	if card.UserId != c.UserId || c.BoardId != card.BoardId || list.UserId != c.UserId {
		return errors.New("board or list does not exist")
	}

	filter := bson.M{"_id": c.ID}
	if _, err := mongodb.Cards.UpdateOne(context.TODO(), filter, bson.M{"$set": c}); err != nil {
		return err
	}
	return nil
}
