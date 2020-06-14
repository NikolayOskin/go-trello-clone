package handlers

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateBoard(b model.Board) error {
	col := mongodb.Client.Database("trello").Collection("boards")
	if _, err := col.InsertOne(context.TODO(), b); err != nil {
		return err
	}
	return nil
}

func UpdateBoard(b model.Board) error {
	col := mongodb.Client.Database("trello").Collection("boards")
	f := bson.M{"_id": b.ID, "user_id": b.UserId}
	if _, err := col.UpdateOne(context.TODO(), f, bson.M{"$set": b}); err != nil {
		return err
	}
	return nil
}
