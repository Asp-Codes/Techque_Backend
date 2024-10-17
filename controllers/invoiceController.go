package controller

import (
	"golang-techque/database"
	"golang-techque/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

// type InvoiceViewFormat struct {
// 	Invoice_id       string
// 	Payment_method   string
// 	Order_id         string
// 	Payment_status   *string
// 	Payment_due      interface{}
// 	Table_number     interface{}
// 	Payment_due_date time.Time
// 	Order_details    interface{}
// }

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var results []bson.M

		result, err := invoiceCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(500, gin.H{"error": "Error while fetching invoices"})
			return
		}

		err = result.All(ctx, &results)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error while decoding invoices"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, results)
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var invoice models.Invoice
		invoiceId := c.Param("invoice_id")

		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching invoice"})
			defer cancel()
			return
		}

		defer cancel()
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
