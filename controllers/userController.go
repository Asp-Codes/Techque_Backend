package controller

import (
	"context"
	"fmt"
	"golang-techque/database"
	helper "golang-techque/helpers"
	"golang-techque/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var users []models.User

		result, err := userCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(500, gin.H{"error": "error fetching users"})
			defer cancel()
			return
		}

		err = result.All(ctx, &users)

		if err != nil {
			c.JSON(500, gin.H{"error": "error decoding users"})
			defer cancel()
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
			defer cancel()
			return
		}

		defer cancel()
		c.JSON(200, user)

	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		validateError := validate.Struct(user)
		if validateError != nil {
			c.JSON(400, gin.H{"error": validateError.Error()})
			defer cancel()
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		if err != nil {
			c.JSON(500, gin.H{"error": "error occured while checking for the user"})
			defer cancel()
			return
		}

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		if err != nil {
			c.JSON(500, gin.H{"error": "error occured while checking for the user"})
			defer cancel()
			return
		}

		if count > 0 {
			c.JSON(400, gin.H{"error": "user already exists ty loginging with different user credentials"})
			defer cancel()
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		result, err := userCollection.InsertOne(ctx, user)

		if err != nil {
			c.JSON(500, gin.H{"error": "error occured while inserting user-user not created"})
			defer cancel()
			return
		}

		defer cancel()
		c.JSON(200, result)

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
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		var foundUser models.User

		err := c.BindJSON(&user)

		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		result := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		if result != nil {
			c.JSON(400, gin.H{"error": "user not found"})
			defer cancel()
			return
		}

		verify, _ := VerifyPassword(*user.Password, *foundUser.Password)

		if !verify {
			c.JSON(400, gin.H{"error": "invalid credentials"})
			defer cancel()
			return
		}

		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)

		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		defer cancel()
		c.JSON(http.StatusOK, foundUser)
		//convert the json data to something golang understands

		// find user in db with email

		// compare password

		// generate token (generateAllTokens function in the helper)

		// update tokens as well (updateAllTokens function in the helper)

		// send response
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
		// return ""
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))

	check := true
	msg := "Password is correct"

	if err != nil {
		msg = fmt.Sprintf("login password failed")
		check = false
	}

	return check, msg
}
