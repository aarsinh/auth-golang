package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName string             `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string             `json:"last_name" validate:"required,min=2,max=100"`
	Email     string             `json:"email" validate:"email,required"`
	Password  string             `json:"password" validate:"required,min=8"`
	UserID    string             `json:"user_id"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=8"`
}
