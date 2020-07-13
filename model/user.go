package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email             string             `json:"email" validate:"required,email"`
	Password          string             `json:"password" validate:"required,passwrd"`
	Verified          bool               `json:"verified"`
	VerificationCode  string             `bson:"confirm_code,omitempty" json:"confirm_code"`
	ResetPasswordCode string             `bson:"reset_password_code,omitempty" json:"reset_password_code"`
}
