package handlers

import (
	"context"
	"errors"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"github.com/NikolayOskin/go-trello-clone/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func HandleCreateList(list model.List) error {
	repo := repository.Boards{}
	board, err := repo.GetById(list.BoardId)
	if err != nil {
		return err
	}
	if board == nil || board.UserId != list.UserId {
		return errors.New("board for this user does not exist")
	}
	col := mongodb.Client.Database("trello").Collection("lists")
	if _, err := col.InsertOne(context.TODO(), list); err != nil {
		return err
	}
	return nil
}

func HandleUpdateList(l model.List) error {
	col := mongodb.Client.Database("trello").Collection("lists")
	f := bson.M{"_id": l.ID, "user_id": l.UserId}
	if _, err := col.UpdateOne(context.TODO(), f, bson.M{"$set": l}); err != nil {
		return err
	}
	return nil
}
