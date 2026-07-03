package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"golang-restaurant-management/database"
	helpers "golang-restaurant-management/helpers"
	"golang-restaurant-management/models"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching users"})
			return
		}

		var users []models.User
		if err = result.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding users"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching user"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		if user.Email != nil {
			emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": *user.Email})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking email"})
				return
			}
			if emailCount > 0 {
				c.JSON(http.StatusConflict, gin.H{"error": "This email already exists"})
				return
			}
		}

		if user.Phone != nil {
			phoneCount, err := userCollection.CountDocuments(ctx, bson.M{"phone": *user.Phone})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking phone"})
				return
			}
			if phoneCount > 0 {
				c.JSON(http.StatusConflict, gin.H{"error": "This phone number already exists"})
				return
			}
		}

		if user.Password != nil {
			password := HashPassword(*user.Password)
			user.Password = &password
		}

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		if user.Email != nil && user.FirstName != nil && user.LastName != nil {
			token, refreshToken, _ := helpers.GenerateAllToken(*user.Email, *user.FirstName, *user.LastName, user.User_id)
			user.Token = &token
			user.RefreshToken = &refreshToken
		}

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while inserting user"})
			return
		}

		c.JSON(http.StatusCreated, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		if user.Email == nil || user.Password == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": *user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*foundUser.Password, *user.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email != nil && foundUser.FirstName != nil && foundUser.LastName != nil {
			token, refreshToken, _ := helpers.GenerateAllToken(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.User_id)
			if err := helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating tokens"})
				return
			}
		}

		c.JSON(http.StatusOK, foundUser)
	}
}

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashedPassword)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	if err != nil {
		return false, "Invalid email or password"
	}
	return true, "Email and password are correct"
}
