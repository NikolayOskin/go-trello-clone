package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Card struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Text     string             `json:"text" bson:"text" validate:"required,max=2000"`
	ListId   string             `json:"list_id" bson:"list_id" validate:"required"`
	BoardId  string             `json:"board_id" bson:"board_id" validate:"required"`
	UserId   string             `json:"user_id" bson:"user_id"` // taken from JWT token
	Position uint               `json:"pos" bson:"pos" validate:"required"`
}
