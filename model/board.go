package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Board struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title  string             `json:"title" bson:"title" validate:"required,max=20"`
	UserId string             `json:"user_id,omitempty" bson:"user_id"` // taken from JWT token
	Lists  []List             `json:"lists"`
}
