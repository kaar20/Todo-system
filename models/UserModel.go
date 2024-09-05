package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserModel struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `json:"name",validate:"required"`
	Phone    string             `json:"phone",validate:"required"`
	Password string             `json:"password",validate:"required"`
}
