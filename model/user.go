package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email                      string             `json:"email" validate:"required,email"`
	Password                   string             `json:"password" validate:"required,min=6,max=80"`
	Verified                   bool               `json:"verified"`
	VerificationCode           string             `bson:"confirm_code,omitempty" json:"confirm_code"`
	ResetPasswordCode          string             `bson:"reset_password_code,omitempty" json:"reset_password_code"`
	ResetPasswordCodeExpiredAt time.Time          `bson:"reset_password_expired_at,omitempty"`
}

type ReadUser struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email    string             `json:"email"`
	Verified bool               `json:"verified"`
}

func (u *User) Verify() {
	u.Verified = true
	u.VerificationCode = ""
}
