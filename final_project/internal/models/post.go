package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId"        json:"userId"`
	Title     string             `bson:"title"         json:"title"`
	Content   string             `bson:"content"       json:"content"`
	UpdatedAt time.Time          `bson:"updatedAt"     json:"updatedAt"`
	CreatedAt time.Time          `bson:"createdAt"     json:"createdAt"`
}
