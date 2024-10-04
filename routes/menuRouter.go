package routes

import (
	controller "golang-techque/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/menus", controller.GetMenuItems())
	incomingRoutes.GET("/menus/:menu_id", controller.GetMenuItem())
	incomingRoutes.POST("/menus", controller.CreateMenuItem())
	incomingRoutes.PATCH("/menus/:menu_id", controller.UpdateMenuItem())
	incomingRoutes.DELETE("/menus/:menu_id", controller.DeleteMenuItem())
}
