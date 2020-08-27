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

type List struct{}

func (h List) CreateList(ctx context.Context, list model.List) (string, error) {
	board, err := repository.Boards.GetById(ctx, list.BoardId)
	if err != nil {
		return "", err
	}
	if board == nil || board.UserId != list.UserId {
		return "", errors.New("board for this user does not exist")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := db.Lists.InsertOne(ctx, list)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (h List) UpdateList(ctx context.Context, l model.UpdateList) error {
	f := bson.M{"_id": l.ID, "user_id": l.UserId}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := db.Lists.UpdateOne(ctx, f, bson.M{"$set": l}); err != nil {
		return err
	}
	return nil
}

func (h List) DeleteList(ctx context.Context, listId primitive.ObjectID, u model.User) error {
	// first deleting cards associated with list, then delete list
	f := bson.M{"list_id": listId.Hex(), "user_id": u.ID.Hex()}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := db.Cards.DeleteMany(ctx, f); err != nil {
		return err
	}
	f = bson.M{"_id": listId, "user_id": u.ID.Hex()}

	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := db.Lists.DeleteOne(ctx, f); err != nil {
		return err
	}
	return nil
}
