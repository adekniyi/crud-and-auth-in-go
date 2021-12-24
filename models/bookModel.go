package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	Id         primitive.ObjectID `bson:"_id"`
	Name       string             `json:"name" validate:"required,min=2,max=100"`
	Genre      string             `json:"genre" validate:"required,min=2,max=100"`
	Created_by string             `json:"created_by"`
	Created_at time.Time          `json:"created_at"`
}
