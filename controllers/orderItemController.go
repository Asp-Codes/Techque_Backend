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

		result, err := orderItemCollection.Find(ctx, bson.M{})
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

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "order_id", Value: orderId}}}} //getting the orders seperated here
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "food"}, {Key: "localField", Value: "food_id"}, {Key: "foreignField", Value: "food_id"}, {Key: "as", Value: "food"}}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$food"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

	lookupOrderStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "order"}, {Key: "localField", Value: "order_id"}, {Key: "foreignField", Value: "order_id"}, {"as", "order"}}}}
	unwindOrderStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$order"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

	lookupTableStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "table"}, {Key: "localField", Value: "table_id"}, {Key: "foreignField", Value: "table_id"}, {Key: "as", Value: "table"}}}}
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

	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "order_id", Value: "$order_id"}, {Key: "table_id", Value: "$table_id"}, {Key: "table_number", Value: "$table_number"}}}, {Key: "payment_due", Value: bson.D{{Key: "$sum", Value: "$amount"}}}, {Key: "table_count", Value: bson.D{{Key: "$sum", Value: 1}}}, {Key: "order_items", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}

	projectStage2 := bson.D{
		{
			Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "payment_due", Value: 1},
				{Key: "table_count", Value: 1},
				{Key: "table_number", Value: "$_id.table_number"},
				{Key: "order_items", Value: 1},
			}},
	}

	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
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

		insertedOrderItems, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
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

		result, err := orderItemCollection.UpdateOne(
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
