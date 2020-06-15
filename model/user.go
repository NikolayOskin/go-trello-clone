package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email            string             `json:"email" validate:"required,email"`
	Password         string             `json:"password" validate:"required,min=6"`
	Verified         bool               `json:"verified"`
	VerificationCode int                `bson:"confirm_code,omitempty" json:"confirm_code"`
}
