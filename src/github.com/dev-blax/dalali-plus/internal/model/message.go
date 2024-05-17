package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Content  string             `bson:"content"`
	Time     int64              `bson:"time"`
}
