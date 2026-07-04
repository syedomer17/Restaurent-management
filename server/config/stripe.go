package config

import (
	"os"

	"github.com/stripe/stripe-go"
)

func InitStripe() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}