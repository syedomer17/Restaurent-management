package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang-restaurant-management/controllers"
)

func OrderRouter(incomeRouters *gin.Engine) {
	incomeRouters.GET("/orders", controller.GetOrders())
	incomeRouters.GET("/orders/:order_id", controller.GetOrder())
	incomeRouters.POST("/orders", controller.CreateOrder())
	incomeRouters.PATCH("/orders/:order_id", controller.UpdateOrder())
	incomeRouters.DELETE("/orders/:order_id", controller.DeleteOrder())
}