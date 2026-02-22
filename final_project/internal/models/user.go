package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username"      json:"username"`
	Password  string             `bson:"password"      json:"-"`
	Email     string             `bson:"email"         json:"email"`
	UpdatedAt time.Time          `bson:"updatedAt"     json:"updatedAt"`
	CreatedAt time.Time          `bson:"createdAt"     json:"createdAt"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	UpdatedAt time.Time          `json:"updatedAt"`
	CreatedAt time.Time          `json:"createdAt"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		UpdatedAt: u.UpdatedAt,
		CreatedAt: u.CreatedAt,
	}
}
