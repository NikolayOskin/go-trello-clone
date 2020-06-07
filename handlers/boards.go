package handlers

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func HandleCreateBoard(board model.Board) error {
	col := mongodb.Client.Database("trello").Collection("boards")
	_, err := col.InsertOne(context.TODO(), board)
	if err != nil {
		return err
	}
	return nil
}

func HandleUpdateBoard(board model.Board) error {
	col := mongodb.Client.Database("trello").Collection("boards")
	filter := bson.M{"_id": board.ID, "user_id": board.UserId}
	_, err := col.UpdateOne(context.TODO(), filter, bson.M{"$set": board})
	if err != nil {
		return err
	}
	return nil
}
