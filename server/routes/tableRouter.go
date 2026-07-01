package routes

import (
	"github.com/gin-gonic/gin"
	 controller "golang-restaurant-management/controllers"
)

func TableRouter(incomeRouters *gin.Engine) {
	incomeRouters.GET("/tables", controller.GetTables())
	incomeRouters.GET("/tables/:table_id", controller.GetTable())
	incomeRouters.POST("/tables", controller.CreateTable())
	incomeRouters.PATCH("/tables/:table_id", controller.UpdateTable())
	incomeRouters.DELETE("/tables/:table_id", controller.DeleteTable())
}