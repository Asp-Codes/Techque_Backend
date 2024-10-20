package controller

import (
	"golang-techque/database"
	"golang-techque/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

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
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

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
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var invoice models.Invoice
		var order models.Order

		err := c.BindJSON(&invoice)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		validationErr := validate.Struct(invoice)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			defer cancel()
			return
		}

		result := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.Order_id}).Decode(&order)

		if result != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching order"})
			defer cancel()
			return
		}

		status := "Pending"
		if invoice.Payment_status == nil {
			invoice.Payment_status = &status
		}

		invoice.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.ID = primitive.NewObjectID()
		invoice.Invoice_id = invoice.ID.Hex()
		// invoice.Payment_due_date = invoice.Created_at.AddDate(0, 0, 1)

		results, insertErr := invoiceCollection.InsertOne(ctx, invoice)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting invoice"})
			defer cancel()
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, results)

	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		invoice_id := c.Param("invoice_id")

		filter := bson.M{"invoice_id": invoice_id}

		upsert := true

		var invoice models.Invoice

		err := c.BindJSON(&invoice)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updatedObj primitive.D

		if invoice.Payment_method != nil {
			updatedObj = append(updatedObj, bson.E{Key: "payment_method", Value: invoice.Payment_method})
		}

		status := "Pending"
		if invoice.Payment_status == nil {
			invoice.Payment_status = &status
		}

		updatedObj = append(updatedObj, bson.E{Key: "payment_status", Value: invoice.Payment_status})

		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updatedObj = append(updatedObj, bson.E{Key: "updated_at", Value: invoice.Updated_at})

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := invoiceCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{Key: "$set", Value: updatedObj}},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating invoice"})
			defer cancel()
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}

func DeleteInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
