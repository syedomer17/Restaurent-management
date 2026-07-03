package controllers

import (
	"context"
	"golang-restaurant-management/database"
	"golang-restaurant-management/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceViewFormatter struct {
	Invoice_id       string
	Payment_method   string
	Order_id         string
	Payment_status   string
	Payment_due      interface{}
	Table_number     interface{}
	Payment_due_date time.Time
	Order_details    interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var invoices models.Invoice 

		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoices.Invoice_id}).Decode(&invoices)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching the invoice"})
			return
		}

		var invoiceView InvoiceViewFormatter

		allOrderItems, err := ItemByOrder(invoices.Order_id)
		invoiceView.Order_id = invoices.Order_id
		invoiceView.Payment_due_date = invoices.Payment_due_date
		
		invoiceView.Payment_method = "null"
		if invoices.Payment_method != nil {
			invoiceView.Payment_method = *invoices.Payment_method
		}

		invoiceView.Invoice_id = invoices.Invoice_id
		invoiceView.Payment_status = *&invoices.Payment_status
		invoiceView.Payment_due = allOrderItems[0]["payment_due"]
		invoiceView.Table_number = allOrderItems[0]["table_number"]
		invoiceView.Order_details = allOrderItems[0]["order_items"]

		c.JSON(http.StatusOK, invoiceView)

	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		invoiceId := c.Param("invoice_id")
		var invoice bson.M

		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching the invoice"})
			return
		}

		c.JSON(http.StatusOK, invoice)
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func DeleteInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
