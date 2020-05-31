package handlers

import (
	"context"
	"fmt"
	"github.com/NikolayOskin/go-trello-clone/model"
	"github.com/NikolayOskin/go-trello-clone/mongodb"
)

func HandleAddBoard(board model.Board) {
	collection := mongodb.Client.Database("trello").Collection("boards")

	insertedBoard, err := collection.InsertOne(context.TODO(), board)
	if err != nil {
		panic(err)
	}
	fmt.Println("inserted board: ", insertedBoard.InsertedID)
}
