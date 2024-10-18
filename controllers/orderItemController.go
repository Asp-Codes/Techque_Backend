package controller

import (
	"context"
	"golang-techque/database"
	"golang-techque/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderItemPack struct {
	Table_id   *string
	OrderItems []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var allorders []models.OrderItem

		result, err := orderItemCollectioon.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(500, gin.H{"error": "error fetching order items"})
			return
		}

		err = result.All(ctx, &allorders)
		if err != nil {
			c.JSON(500, gin.H{"error": "error decoding order items"})
			return
		}
		defer cancel()
		c.JSON(200, allorders)

	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")

		allOrderItems, err := ItemsByOrder(orderId)
		if err != nil {
			c.JSON(500, gin.H{"error": "error fetching order items"})
			return
		}
		c.JSON(200, allOrderItems)
	}
}

func ItemsByOrder(orderId string) (orderItems []primitive.M, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	//using pipeline

	matchStage := bson.D{{"$match", bson.D{{"order_id", orderId}}}} //getting the orders seperated here
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localField", "food_id"}, {"foreignField", "food_id"}, {"as", "food"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupOrderStage := bson.D{{"$lookup", bson.D{{"from", "order"}, {"localField", "order_id"}, {"foreignField", "order_id"}, {"as", "order"}}}}
	unwindOrderStage := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupTableStage := bson.D{{"$lookup", bson.D{{"from", "table"}, {"localField", "table_id"}, {"foreignField", "table_id"}, {"as", "table"}}}}
	unwindTableStage := bson.D{{"$unwind", bson.D{{"path", "$table"}, {"preserveNullAndEmptyArrays", true}}}}

	projectStage := bson.D{
		{
			"$project", bson.D{
				{"id", 0},
				{"amount", "$food.price"},
				{"total_count", 1},
				{"food_name", "$food.name"},
				{"food_image", "$food.food_image"},
				{"table_number", "$table.table_number"},
				{"table_id", "$table.table_id"},
				{"order_id", "$order.order_id"},
				{"price", "$food.price"},
				{"quantity", 1},
			}},
	}

	groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"order_id", "$order_id"}, {"table_id", "$table_id"}, {"table_number", "$table_number"}}}, {"payment_due", bson.D{{"$sum", "$amount"}}}, {"table_count", bson.D{{"$sum", 1}}}, {"order_items", bson.D{{"$push", bson.D{{"food_name", "$food_name"}, {"food_image", "$food_image"}, {"price", "$price"}, {"quantity", "$quantity"}}}}}}}}

	projectStage2 := bson.D{
		{
			"$project", bson.D{
				{"_id", 0},
				{"payment_due", 1},
				{"table_count", 1},
				{"table_number", "$_id.table_number"},
				{"order_items", 1},
			}},
	}

	result, err := orderItemCollection.Aggregrate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	})

	if err != nil {
		return nil, err
	}

	err = result.All(ctx, &orderItems)
	if err != nil {
		return nil, err
	}

	defer cancel()

	return orderItems, err

}
func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		orderItemId := c.Param("order_item_id")

		var orderItem models.OrderItem

		err := orderItemCollection.FindOne(ctx, bson.M{"order_item": orderItemId}).Decode(&orderItem)

		if err != nil {
			c.JSON(500, gin.H{"error": "error fetching order item"})
			return
		}
		defer cancel()
		c.JSON(200, orderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var order models.Order
		var orderItemPack orderItemPack

		err := c.BindJSON(&orderItemPack)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.Table_id = orderItemPack.Table_id
		order_id := OrderItemOrderCreator(order)
		orderItemsToBeInserted := []interface{}{}

		for _, orderItem := range orderItemPack.OrderItems {
			orderItem.Order_id = order_id

			validationErr := validate.Struct(orderItem)
			if validationErr != nil {
				c.JSON(400, gin.H{"error": validationErr.Error()})
				return
			}

			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.ID = primitive.NewObjectID()
			orderItem.Order_item_id = orderItem.ID.Hex()

			var num = toFixed(*orderItem.Unit_Price, 2)
			orderItem.Unit_Price = &num

			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}

		insertedOrderItems, err := orderItemCollectioon.InsertMany(ctx, orderItemsToBeInserted)
		if err != nil {
			c.JSON(500, gin.H{"error": "error inserting order items"})
			return
		}
		defer cancel()
		c.JSON(200, insertedOrderItems)

	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		OrderItemId := c.Param("order_item_id")

		var orderItem models.OrderItem
		filter := bson.M{"order_item_id": OrderItemId}

		var updatedObj primitive.D

		if orderItem.Unit_Price != nil {
			updatedObj = append(updatedObj, bson.E{Key: "order_item_id", Value: *&orderItem.Order_item_id})
		}

		if orderItem.Quantity != nil {
			updatedObj = append(updatedObj, bson.E{Key: "quantity", Value: *orderItem.Quantity})
		}

		if orderItem.Food_id != nil {
			updatedObj = append(updatedObj, bson.E{Key: "food_id", Value: *orderItem.Food_id})
		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updatedObj = append(updatedObj, bson.E{Key: "updated_at", Value: orderItem.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollectioon.UpdateOne(
			ctx,
			filter,
			bson.D{{
				Key: "$set", Value: updatedObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(500, gin.H{"error": "error updating order item"})
			return
		}

		defer cancel()
		c.JSON(200, result)

	}
}

func DeleteOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
