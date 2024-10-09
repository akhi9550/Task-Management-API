package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `json:"Name"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
}

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      string             `json:"user_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"created_at"`
}
