package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	OrderID primitive.ObjectID `bson:"order_id" json:"order_id" validate:"required"`

	UserID primitive.ObjectID `bson:"user_id" json:"user_id" validate:"required"`

	PaymentIntentID string `bson:"payment_intent_id" json:"payment_intent_id"`

	CustomerID string `bson:"customer_id,omitempty" json:"customer_id,omitempty"`

	Amount int64 `bson:"amount" json:"amount"`

	Currency string `bson:"currency" json:"currency"`

	Status string `bson:"status" json:"status"`

	PaymentMethod string `bson:"payment_method,omitempty" json:"payment_method,omitempty"`

	ReceiptURL string `bson:"receipt_url,omitempty" json:"receipt_url,omitempty"`

	RefundID string `bson:"refund_id,omitempty" json:"refund_id,omitempty"`

	FailureReason string `bson:"failure_reason,omitempty" json:"failure_reason,omitempty"`

	CreatedAt time.Time `bson:"created_at"`

	UpdatedAt time.Time `bson:"updated_at"`
}