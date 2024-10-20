package routes

import (
	controller "golang-techque/controllers"

	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/orderItems", controller.GetOrderItems())
	incomingRoutes.GET("/orderItems/:orderItem_id", controller.GetOrderItem())           //error
	incomingRoutes.GET("/orderItems/order/:order_id", controller.GetOrderItemsByOrder()) //error here as well
	incomingRoutes.POST("/orderItems", controller.CreateOrderItem())
	incomingRoutes.PATCH("/orderItems/:orderItem_id", controller.UpdateOrderItem())
	incomingRoutes.DELETE("/orderItems/:orderItem_id", controller.DeleteOrderItem())
}

//some signifcant issue here!
