package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Order_Date  time.Time          `json:"order_date" validate:"required"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Order_id    string             `json:"order_id" validate:"required"`
	Table_id    string             `json:"table_id" validate:"required"`
}