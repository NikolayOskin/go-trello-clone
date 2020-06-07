package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Board struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title  string             `json:"title" bson:"title"`
	UserId string             `json:"user_id,omitempty" bson:"user_id"`
}
