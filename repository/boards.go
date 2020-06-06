package repository

import (
	"context"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type Boards struct{}

func (b *Boards) FetchByUser(user model.User) []model.Board {
	var boards []model.Board

	collection := mongodb.Client.Database("trello").Collection("boards")
	filter := bson.D{{"user_id", user.ID.Hex()}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &boards); err != nil {
		log.Fatal(err)
	}

	return boards
}

func (b *Boards) GetById(id string) (*model.Board, error) {
	var board model.Board

	collection := mongodb.Client.Database("trello").Collection("boards")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if err := collection.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&board); err != nil {
		return nil, err
	}
	return &board, nil
}
