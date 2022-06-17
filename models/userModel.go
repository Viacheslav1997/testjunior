package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id            primitive.ObjectID `bson:"_id"`
	Refresh_token *string            `json:"refresh_token"`
	User_id       string             `json:"user_id"`
}
