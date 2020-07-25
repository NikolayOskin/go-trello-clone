package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type List struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title    string             `json:"title" bson:"title" validate:"required,max=20"`
	BoardId  string             `json:"board_id" bson:"board_id" validate:"required"`
	UserId   string             `json:"user_id" bson:"user_id"` // taken from JWT token
	Position uint               `json:"pos" bson:"pos" validate:"required"`
	Cards    []Card             `json:"cards"`
}

type UpdateList struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title  string             `json:"title" bson:"title" validate:"required,max=20"`
	UserId string             `json:"user_id" bson:"user_id"` // taken from JWT token
}
