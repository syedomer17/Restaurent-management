package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang-restaurant-management/controllers"
)

func InvoiceRouter(incomeRouters *gin.Engine) {
	incomeRouters.GET("/invoices", controller.GetInvoices())
	incomeRouters.GET("/invoices/:invoice_id", controller.GetInvoice())
	incomeRouters.POST("/invoices", controller.CreateInvoice())
	incomeRouters.PATCH("/invoices/:invoice_id", controller.UpdateInvoice())
	incomeRouters.DELETE("/invoices/:invoice_id", controller.DeleteInvoice())
}