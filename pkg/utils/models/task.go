package models

import "time"

type CreateTask struct {
	Title       string `json:"title" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"required,min=1,max=1000"`
}

type TaskDetails struct {
	ID          string    `bson:"_id"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	CreatedAt   time.Time `bson:"created_at"`
}
