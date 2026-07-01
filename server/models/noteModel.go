package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Text      string             `json:"text" validate:"required"`
	Title     string             `json:"title" validate:"required"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Note_id   string             `json:"note_id"`
}
