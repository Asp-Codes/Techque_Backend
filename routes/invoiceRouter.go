package routes

import (
	controller "golang-techque/controllers"

	"github.com/gin-gonic/gin"
)

func InvoicesRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/invoices", controller.GetInvoices())
	incomingRoutes.GET("/invoices/:invoice_id", controller.GetInvoice())
	incomingRoutes.POST("/invoices", controller.CreateInvoice())
	incomingRoutes.PATCH("/invoices/:invoice_id", controller.UpdateInvoice())
	incomingRoutes.DELETE("/invoices/:invoice_id", controller.DeleteInvoice())
}

//done
