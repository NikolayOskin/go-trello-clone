package repository

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Boards struct{}

func (b *Boards) FetchByUser(user model.User) ([]model.Board, error) {
	var boards []model.Board
	col := mongodb.Client.Database("trello").Collection("boards")
	filter := bson.D{{"user_id", user.ID.Hex()}}
	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &boards); err != nil {
		return nil, err
	}

	return boards, nil
}

func (b *Boards) GetById(id string) (*model.Board, error) {
	var board model.Board
	col := mongodb.Client.Database("trello").Collection("boards")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := col.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&board); err != nil {
		return nil, err
	}

	return &board, nil
}
