package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var Client *mongo.Client
var Users *mongo.Collection
var Boards *mongo.Collection
var Lists *mongo.Collection
var Cards *mongo.Collection

func InitDB() {
	login := os.Getenv("MONGODB_LOGIN")
	pass := os.Getenv("MONGODB_PASSWORD")
	host := os.Getenv("MONGODB_HOST")
	port := os.Getenv("MONGODB_PORT")
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%v:%v@%v:%v", login, pass, host, port))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	Client = client
	Users = client.Database("trello").Collection("users")
	Boards = client.Database("trello").Collection("boards")
	Lists = client.Database("trello").Collection("lists")
	Cards = client.Database("trello").Collection("cards")
}
