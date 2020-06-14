package handlers

import (
	"context"
	"errors"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"github.com/NikolayOskin/go-trello-clone/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateCard(c model.Card) error {
	repo := repository.Lists{}
	list, err := repo.GetById(c.ListId)
	if err != nil {
		return err
	}
	if list == nil || list.UserId != c.UserId {
		return errors.New("list for your user does not exist")
	}
	c.BoardId = list.BoardId
	col := mongodb.Client.Database("trello").Collection("cards")
	if _, err = col.InsertOne(context.TODO(), c); err != nil {
		return err
	}
	return nil
}

func UpdateCard(c model.Card) error {
	col := mongodb.Client.Database("trello").Collection("cards")
	filter := bson.M{"_id": c.ID, "user_id": c.UserId}
	if _, err := col.UpdateOne(context.TODO(), filter, bson.M{"$set": c}); err != nil {
		return err
	}
	return nil
}
