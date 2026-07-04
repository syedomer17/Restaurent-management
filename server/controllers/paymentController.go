package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/webhook"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"golang-restaurant-management/database"
	"golang-restaurant-management/models"
)

var PaymentController *mongo.Collection = database.OpenCollection(database.Client, "payments")

type createPaymentIntentRequest struct {
	OrderID            string   `json:"order_id" binding:"required"`
	UserID             string   `json:"user_id"`
	Amount             int64    `json:"amount" binding:"required"`
	Currency           string   `json:"currency"`
	PaymentMethodTypes []string `json:"payment_method_types"`
}

func GetStripeConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"publishable_key": os.Getenv("STRIPE_PUBLISHABLE_KEY"),
			"mode":            "test",
		})
	}
}

func CreatePaymentIntent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createPaymentIntentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Amount <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be greater than zero"})
			return
		}

		currency := strings.ToLower(strings.TrimSpace(req.Currency))
		if currency == "" {
			currency = "inr"
		}

		paymentMethodTypes := normalizePaymentMethodTypes(req.PaymentMethodTypes)
		if len(paymentMethodTypes) == 0 {
			paymentMethodTypes = []string{"card", "upi"}
		}

		secretKey := os.Getenv("STRIPE_SECRET_KEY")
		if secretKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Stripe secret key is not configured"})
			return
		}

		stripe.Key = secretKey

		params := &stripe.PaymentIntentParams{
			Amount:             stripe.Int64(req.Amount),
			Currency:           stripe.String(currency),
			Confirm:            stripe.Bool(false),
			UseStripeSDK:       stripe.Bool(true),
			PaymentMethodTypes: stripe.StringSlice(paymentMethodTypes),
			Metadata: map[string]string{
				"order_id": req.OrderID,
				"user_id":  req.UserID,
			},
			Description: stripe.String("Restaurant order payment"),
		}

		intent, err := paymentintent.New(params)
		if err != nil {
			fallbackParams := &stripe.PaymentIntentParams{
				Amount:             stripe.Int64(req.Amount),
				Currency:           stripe.String(currency),
				Confirm:            stripe.Bool(false),
				UseStripeSDK:       stripe.Bool(true),
				PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
				Metadata: map[string]string{
					"order_id": req.OrderID,
					"user_id":  req.UserID,
				},
				Description: stripe.String("Restaurant order payment"),
			}
			intent, err = paymentintent.New(fallbackParams)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create Stripe payment intent: " + err.Error()})
				return
			}
		}

		var orderID primitive.ObjectID
		if parsedOrderID, err := primitive.ObjectIDFromHex(req.OrderID); err == nil {
			orderID = parsedOrderID
		}

		var userID primitive.ObjectID
		if req.UserID != "" {
			if parsedUserID, err := primitive.ObjectIDFromHex(req.UserID); err == nil {
				userID = parsedUserID
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		paymentRecord := models.Payment{
			ID:              primitive.NewObjectID(),
			OrderID:         orderID,
			UserID:          userID,
			PaymentIntentID: intent.ID,
			Amount:          req.Amount,
			Currency:        strings.ToUpper(currency),
			Status:          string(intent.Status),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			PaymentMethod:   "",
			ReceiptURL:      "",
			FailureReason:   "",
		}

		_, err = PaymentController.InsertOne(ctx, paymentRecord)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "payment intent created but failed to save payment record"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":         "Payment intent created successfully",
			"payment_id":      paymentRecord.ID.Hex(),
			"payment_intent":  intent.ID,
			"client_secret":   intent.ClientSecret,
			"status":          intent.Status,
			"currency":        strings.ToUpper(currency),
			"amount":          req.Amount,
			"payment_methods": paymentMethodTypes,
		})
	}
}

func StripeWebhook() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read request body"})
			return
		}

		signature := c.GetHeader("Stripe-Signature")
		endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
		if endpointSecret == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stripe webhook secret is not configured"})
			return
		}

		event, err := webhook.ConstructEvent(payload, signature, endpointSecret)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Stripe webhook signature"})
			return
		}

		switch event.Type {
		case "payment_intent.succeeded":
			var paymentIntent stripe.PaymentIntent
			if dataBytes, err := json.Marshal(event.Data.Object); err == nil {
				if err := json.Unmarshal(dataBytes, &paymentIntent); err == nil {
					ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
					defer cancel()

					update := bson.M{
						"status":     string(paymentIntent.Status),
						"updated_at": time.Now(),
					}
					if paymentIntent.LatestCharge != nil {
						update["payment_method"] = "card"
					}

					_, _ = PaymentController.UpdateOne(ctx, bson.M{"payment_intent_id": paymentIntent.ID}, bson.D{{Key: "$set", Value: update}})
				}
			}
		case "payment_intent.payment_failed":
			var paymentIntent stripe.PaymentIntent
			if dataBytes, err := json.Marshal(event.Data.Object); err == nil {
				if err := json.Unmarshal(dataBytes, &paymentIntent); err == nil {
					ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
					defer cancel()

					failureReason := ""
					if paymentIntent.LastPaymentError != nil {
						failureReason = paymentIntent.LastPaymentError.Error()
					}

					_, _ = PaymentController.UpdateOne(ctx, bson.M{"payment_intent_id": paymentIntent.ID}, bson.D{{Key: "$set", Value: bson.M{
						"status":         string(paymentIntent.Status),
						"failure_reason": failureReason,
						"updated_at":     time.Now(),
					}}})
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"received": true})
	}
}

func GetPayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		paymentID := c.Param("payment_id")
		var payment models.Payment
		filter := bson.M{"payment_intent_id": paymentID}

		if parsedID, err := primitive.ObjectIDFromHex(paymentID); err == nil {
			filter = bson.M{"$or": []bson.M{
				{"_id": parsedID},
				{"payment_intent_id": paymentID},
			}}
		}

		err := PaymentController.FindOne(ctx, filter).Decode(&payment)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}

		c.JSON(http.StatusOK, payment)
	}
}

func GetPaymentByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderID := c.Param("order_id")
		var payment models.Payment
		filter := bson.M{"order_id": orderID}

		if parsedOrderID, err := primitive.ObjectIDFromHex(orderID); err == nil {
			filter = bson.M{"order_id": parsedOrderID}
		}

		err := PaymentController.FindOne(ctx, filter).Decode(&payment)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found for this order"})
			return
		}

		c.JSON(http.StatusOK, payment)
	}
}

func normalizePaymentMethodTypes(requestedTypes []string) []string {
	if len(requestedTypes) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	normalized := make([]string, 0, len(requestedTypes))
	for _, method := range requestedTypes {
		switch strings.ToLower(strings.TrimSpace(method)) {
		case "card", "credit_card", "debit_card":
			method = "card"
		case "upi", "upi_collect", "upi_qr":
			method = "upi"
		case "link":
			method = "link"
		default:
			method = strings.ToLower(strings.TrimSpace(method))
		}
		if method != "" && !seen[method] {
			seen[method] = true
			normalized = append(normalized, method)
		}
	}

	return normalized
}
