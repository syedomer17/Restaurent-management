package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang-restaurant-management/controllers"
)

func OrderItemRouter(incomeRouters *gin.Engine) {
	incomeRouters.GET("/orderitems", controller.GetOrderItems())
	incomeRouters.GET("/orderitems/:orderitem_id", controller.GetOrderItem())
	incomeRouters.GET("/orderitems-order/:order_id", controller.GetOrderItemsByOrder())
	incomeRouters.POST("/orderitems", controller.CreateOrderItem())
	incomeRouters.PATCH("/orderitems/:orderitem_id", controller.UpdateOrderItem())
	incomeRouters.DELETE("/orderitems/:orderitem_id", controller.DeleteOrderItem())
}