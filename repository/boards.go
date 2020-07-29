package repository

import (
	"context"
	"time"

	"github.com/NikolayOskin/go-trello-clone/db"
	"github.com/NikolayOskin/go-trello-clone/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Boards struct{}

func (b *Boards) FetchByUser(ctx context.Context, user model.User) ([]model.Board, error) {
	var boards []model.Board
	filter := bson.D{{"user_id", user.ID.Hex()}}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := db.Boards.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &boards); err != nil {
		return nil, err
	}

	return boards, nil
}

func (b *Boards) GetById(ctx context.Context, id string) (*model.Board, error) {
	var board model.Board
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Boards.FindOne(ctx, bson.M{"_id": objId}).Decode(&board); err != nil {
		return nil, err
	}

	return &board, nil
}
