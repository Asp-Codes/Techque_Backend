package controller

import (
	"context"
	"golang-techque/database"
	"golang-techque/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement logic to fetch all tables from the database
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var allTables []bson.M
		result, err := tableCollection.Find(context.TODO(), bson.M{})

		if err != nil {
			c.JSON(500, gin.H{"error": "error fetching tables"})
			return
		}

		err = result.All(ctx, &allTables)
		if err != nil {
			log.Fatal(err)
		}
		defer cancel()
		c.JSON(http.StatusOK, allTables)
	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement logic to fetch a table by its id from the database
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		tableId := c.Param("table_id")

		var table models.Table

		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching table"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, table)

	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement logic to create a new table in the database
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var table models.Table

		err := c.BindJSON(&table)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		validationErr := validate.Struct(table)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			defer cancel()
			return
		}

		table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()

		result, insertErr := tableCollection.InsertOne(ctx, table)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error inserting table"})
			defer cancel()
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement logic to update an existing table in the database
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		err := c.BindJSON(&table)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tableId := c.Param("table_id")
		filter := bson.M{"table_id": tableId}
		upsert := true

		var updatedObj primitive.D

		if table.Number_of_guests != nil {
			updatedObj = append(updatedObj, bson.E{Key: "number_of_guests", Value: table.Number_of_guests})
		}

		Status := "Vacant"
		if table.Status == nil {
			table.Status = &Status
		}

		updatedObj = append(updatedObj, bson.E{Key: "table_status", Value: table.Status})

		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updatedObj = append(updatedObj, bson.E{Key: "updated_at", Value: table.Updated_at})

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := tableCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{Key: "$set", Value: updatedObj}},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating table"})
			defer cancel()
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}

func DeleteTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement logic to delete a table from the database
	}
}
