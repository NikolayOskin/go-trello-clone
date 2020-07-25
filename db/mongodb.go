package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
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
	dbname := os.Getenv("MONGODB_DBNAME")

	if os.Getenv("APP_ENV") == "test" {
		dbname = os.Getenv("MONGODB_TEST_DBNAME")
	}

	dbUri := fmt.Sprintf("mongodb://%v:%v@%v:%v", login, pass, host, port)

	clientOptions := options.Client().ApplyURI(dbUri)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(dbname)

	Client = client

	Users = db.Collection("users")
	Boards = db.Collection("boards")
	Lists = db.Collection("lists")
	Cards = db.Collection("cards")

	createIndexes()
}

func createIndexes() {
	createIndex(Users, "email", true)
	createIndex(Boards, "user_id", false)
	createIndex(Lists, "board_id", false)
	createIndex(Cards, "board_id", false)
}

// createIndex - creates an index for a specific field in a collection
func createIndex(collection *mongo.Collection, field string, unique bool) bool {
	mod := mongo.IndexModel{
		Keys:    bson.M{field: 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(unique),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

// FreshDb - delete all data from db and recreate indexes
func FreshDb() {
	if os.Getenv("APP_ENV") != "test" {
		log.Fatal("you can refresh database only in test environment")
	}

	if err := Users.Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}
	if err := Boards.Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}
	if err := Lists.Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}
	if err := Cards.Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}

	createIndexes()
}

// IsDuplicated - check if error is unique index duplicated error
func IsDuplicated(err error) bool {
	if we, ok := err.(mongo.WriteException); ok {
		for _, e := range we.WriteErrors {
			if e.Code == 11000 {
				return true
			}
		}
	}
	return false
}
