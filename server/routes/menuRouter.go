package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang-restaurant-management/controllers"
)

func MenuRouter(incomeRouters *gin.Engine) {
	incomeRouters.GET("/menus", controller.GetMenus())
	incomeRouters.GET("/menus/:menu_id", controller.GetMenu())
	incomeRouters.POST("/menus", controller.CreateMenu())
	incomeRouters.PATCH("/menus/:menu_id", controller.UpdateMenu())
	incomeRouters.DELETE("/menus/:menu_id", controller.DeleteMenu())
}