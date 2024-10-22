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

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

// var validate = validator.New()

func GetMenuItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()

		if err != nil {
			c.JSON(500, gin.H{"error": "error fetching menu items"})
			return
		}

		var allMenus []bson.M //using a slice here of BJSON

		err = result.All(ctx, &allMenus)

		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenus) //BJSON is converted into JSON

	}
}

func GetMenuItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		menuID := c.Param("menu_id")

		var menuItem models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuID}).Decode(&menuItem)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching menu item"})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, menuItem)
	}
}

func CreateMenuItem() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var menuItem models.Menu

		err := c.BindJSON(&menuItem)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(menuItem)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		defer cancel()
		menuItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menuItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menuItem.ID = primitive.NewObjectID()
		menuItem.Menu_id = menuItem.ID.Hex()

		result, err := menuCollection.InsertOne(ctx, menuItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error inserting menu item"})
			return
		}
		c.JSON(http.StatusOK, result)

	}
}

func GetMenuItemsByMenuId() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		menuID := c.Param("menu_id")

		defer cancel()

		matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "menu_id", Value: menuID}}}}

		lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "food"}, {Key: "localField", Value: "menu_id"}, {Key: "foreignField", Value: "menu_id"}, {Key: "as", Value: "menu_items"}}}}

		unWindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$menu_items"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

		projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "category", Value: 1}, {Key: "name", Value: "$menu_items.name"}, {Key: "food_id", Value: "$menu_items.food_id"}, {Key: "price", Value: "$menu_items.price"}, {Key: "food_image", Value: "$menu_items.food_image"}}}}

		cursor, err := menuCollection.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, unWindStage, projectStage})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching menu items"})
			return
		}

		var menuItems []bson.M

		if err = cursor.All(ctx, &menuItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching menu items"})
			return
		}

		c.JSON(http.StatusOK, menuItems)

	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.After(time.Now()) && end.After(time.Now())
}

func UpdateMenuItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		menuID := c.Param("menu_id")
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"menu_id": menuID}

		var updateObj primitive.D

		if menu.Start_date != nil && menu.End_date != nil {
			if !inTimeSpan(*menu.Start_date, *menu.End_date, time.Now()) {
				msg := "kindly recheck the time"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				defer cancel()
				return
			}

			updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.Start_date})
			updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.End_date})

			if menu.Name != "" {
				updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
			}

			if menu.Category != "" {
				updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
			}

			menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

			updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.Updated_at})

			upsert := true
			opt := options.UpdateOptions{
				Upsert: &upsert,
			}

			result, err := menuCollection.UpdateOne(
				ctx,
				filter,
				bson.D{{Key: "$set", Value: updateObj}},
				&opt,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating menu item"})
				defer cancel()
				return
			}
			defer cancel()
			c.JSON(http.StatusOK, result)
		}
	}
}

func DeleteMenuItem() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel() // Ensure the context is canceled at the end

        menuID := c.Param("menu_id") // Get menu_id from URL parameters

        // Attempt to delete the menu item from the database
        result, err := menuCollection.DeleteOne(ctx, bson.M{"menu_id": menuID})
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting menu item"})
            return
        }

        // Check if any document was deleted
        if result.DeletedCount == 0 {
            c.JSON(http.StatusNotFound, gin.H{"error": "menu item not found"})
            return
        }

        // Respond with success
        c.JSON(http.StatusOK, gin.H{"message": "menu item deleted successfully"})
    }
}
