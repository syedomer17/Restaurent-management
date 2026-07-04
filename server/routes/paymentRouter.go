package routes

import (
	"golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func PaymentRoutes(router *gin.Engine) {

	router.GET("/payments/config",
		controllers.GetStripeConfig())

	router.POST("/payments/create-payment-intent",
		controllers.CreatePaymentIntent())

	router.POST("/payments/webhook",
		controllers.StripeWebhook())

	router.GET("/payments/:payment_id",
		controllers.GetPayment())

	router.GET("/payments/order/:order_id",
		controllers.GetPaymentByOrder())
}
