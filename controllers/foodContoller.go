package controller

import (
	"context"
	"golang-techque/database"
	"golang-techque/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("food_id")
		var food models.Food

		foodCollection.FindOne()
	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func round(num float64) int {
	return 0
}

func toFixed(num float64, precision int) float64 {
	return 0
}
func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func DeleteFood() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
