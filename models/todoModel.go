package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoModel struct {
	ID           primitive.ObjectID `bson:"_id"`
	Todo_id      string             `json:"todo"`
	Title        string             `json:"title",validate:"required"`
	Description  string             `json:"description",validate:"required"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
	Is_completed bool               `json:"is_completed"`
	User         string             `json:"user" ,validate:"required"`
}
