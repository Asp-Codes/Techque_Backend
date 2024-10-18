package controller

import (
	"context"
	"golang-techque/database"
	"golang-techque/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var users []models.User

		result, err := userCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(500, gin.H{"error": "error fetching users"})
			return
		}

		err = result.All(ctx, &users)

		if err != nil {
			c.JSON(500, gin.H{"error": "error decoding users"})
			return
		}

		defer cancel()
		c.JSON(200, users)

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		userId := c.Param("user_id")

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		if err != nil {
			c.JSON(500, gin.H{"error": "error fetching user"})
			return
		}

		defer cancel()
		c.JSON(200, user)

	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		// convert the JSON data to smthng golanng understand
		// validate data
		// check if user exists
		// check both email and password
		// create extra details like upadated at and created at
		// hash password
		// save to db
		// generate token (generateAllTokens function in the helper)
		// send response

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		//convert the json data to something golang understands

		// find user in db with email

		// compare password

		// generate token (generateAllTokens function in the helper)

		// update tokens as well (updateAllTokens function in the helper)

		// send response
	}
}

func HashPassword(password string) string {
	return ""
}

func VerifyPassword(userPassword string, providedPassword string) bool {
	return false
}
