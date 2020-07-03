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

func CreateList(list model.List) (string, error) {
	repo := repository.Boards{}
	board, err := repo.GetById(list.BoardId)
	if err != nil {
		return "", err
	}
	if board == nil || board.UserId != list.UserId {
		return "", errors.New("board for this user does not exist")
	}
	res, err := mongodb.Lists.InsertOne(context.TODO(), list)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func UpdateList(l model.List) error {
	f := bson.M{"_id": l.ID, "user_id": l.UserId}
	if _, err := mongodb.Lists.UpdateOne(context.TODO(), f, bson.M{"$set": l}); err != nil {
		return err
	}
	return nil
}

func DeleteList(listId primitive.ObjectID, u model.User) error {
	// first deleting cards associated with list, then delete list
	f := bson.M{"list_id": listId.Hex(), "user_id": u.ID.Hex()}
	if _, err := mongodb.Cards.DeleteMany(context.TODO(), f); err != nil {
		return err
	}
	f = bson.M{"_id": listId, "user_id": u.ID.Hex()}
	if _, err := mongodb.Lists.DeleteOne(context.TODO(), f); err != nil {
		return err
	}
	return nil
}
