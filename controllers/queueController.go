package controller

import (
	"context"
	"fmt"
	"golang-techque/database"
	"golang-techque/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var queueCollection *mongo.Collection = database.OpenCollection(database.Client, "queue")

// GetQueue is a function to get all queues
func GetQueues() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var results []models.Queue

		cursor, err := queueCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(500, gin.H{"error": "Error while fetching queues"})
			defer cancel()
			return
		}

		err = cursor.All(ctx, &results)

		if err != nil {
			c.JSON(500, gin.H{"error": "Error while decoding queues"})
			defer cancel()
			return
		}

		defer cancel()
		c.JSON(200, results)

	}
}

// GetQueue is a function to get a queue by its id
func GetQueue() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		queueId := c.Param("queue_id")
		fmt.Println(queueId)

		var result models.Queue

		err := queueCollection.FindOne(ctx, bson.M{"queue_id": queueId}).Decode(&result)

		if err == mongo.ErrNoDocuments {
			c.JSON(404, gin.H{"error": "Queue not found"})
			defer cancel()
			return
		} else if err != nil {
			c.JSON(500, gin.H{"error": "Error while fetching queue"})
			defer cancel()
			return
		}

		defer cancel()
		c.JSON(200, result)

	}
}

func UpdateQueue() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		queueId := c.Param("queue_id")

		if queueId == "" {
			c.JSON(400, gin.H{"error": "Queue ID is required"})
			defer cancel()
			return
		}

		var updatedObj primitive.D
		var queue models.Queue

		var err = queueCollection.FindOne(ctx, bson.M{"queue_id": queueId}).Decode(&queue)

		if err != nil {
			c.JSON(404, gin.H{"error": "Queue not found"})
			defer cancel()
			return
		}

		// if queue.Current_order != -1 && queue.Total_orders != -1 {
		// 	queue.Current_order = queue.Current_order + 1
		// 	updatedObj = append(updatedObj, bson.E{Key: "current_order", Value: queue.Current_order})
		// }

		// err := queueCollection.Find(ctx, bson.M{"queue_id": queue.Queue_id}).Decode(&queue)

		// if queue.Total_orders != -1 {
		queue.Total_orders = queue.Total_orders + 1
		updatedObj = append(updatedObj, bson.E{Key: "total_orders", Value: queue.Total_orders})
		// }

		queue.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		updatedObj = append(updatedObj, bson.E{Key: "updated_at", Value: queue.Updated_at})

		upsert := true
		filter := bson.M{"queue_id": queue.Queue_id}
		opts := options.UpdateOptions{
			Upsert: &upsert,
		}

		_, err = queueCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updatedObj}}, &opts)

		if err != nil {
			c.JSON(500, gin.H{"error": "Error while updating queue"})
			defer cancel()
			return
		}

		defer cancel()
		c.JSON(200, queue)

	}
}

func CreateQueue() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var queue models.Queue

		if err := c.BindJSON(&queue); err != nil {
			c.JSON(400, gin.H{"error": "Unable to bind data"})
			defer cancel()
			return
		}

		queue.ID = primitive.NewObjectID()
		queue.Queue_id = (queue.ID).Hex()
		queue.Total_orders = 1
		queue.Current_order = 1
		queue.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		queue.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		_, err := queueCollection.InsertOne(ctx, queue)

		if err != nil {
			c.JSON(500, gin.H{"error": "Error while creating queue"})
			defer cancel()
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, queue)

	}
}

func UpdateQueueOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		queueId := c.Param("queue_id")

		if queueId == "" {
			c.JSON(400, gin.H{"error": "Queue ID is required"})
			defer cancel()
			return
		}

		var updatedObj primitive.D
		var queue models.Queue

		var err = queueCollection.FindOne(ctx, bson.M{"queue_id": queueId}).Decode(&queue)

		if err != nil {
			c.JSON(404, gin.H{"error": "Queue not found"})
			defer cancel()
			return
		}

		if queue.Current_order != -1 && queue.Total_orders > queue.Current_order {
			queue.Current_order = queue.Current_order + 1
			updatedObj = append(updatedObj, bson.E{Key: "current_order", Value: queue.Current_order})
		} else if queue.Total_orders <= queue.Current_order {
			queue.Current_order = 0
			queue.Total_orders = 0
			fmt.Println("No active Orders Buddy")
			updatedObj = append(updatedObj, bson.E{Key: "current_order", Value: queue.Current_order})
			updatedObj = append(updatedObj, bson.E{Key: "total_orders", Value: queue.Total_orders})
		}

		queue.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		updatedObj = append(updatedObj, bson.E{Key: "updated_at", Value: queue.Updated_at})

		upsert := true
		filter := bson.M{"queue_id": queue.Queue_id}
		opts := options.UpdateOptions{
			Upsert: &upsert,
		}

		_, err = queueCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updatedObj}}, &opts)

		if err != nil {
			c.JSON(500, gin.H{"error": "Error while updating queue"})
			defer cancel()
			return
		}

		defer cancel()
		c.JSON(200, queue)
	}
}
