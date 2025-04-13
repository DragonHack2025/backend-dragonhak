package handlers

import (
	"context"
	"net/http"
	"os"
	"time"

	"backend-dragonhak/auth"
	"backend-dragonhak/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetUsers handles getting all users
func GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.User
	cursor, err := Collections.Users.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

// GetUser handles getting a single user
func GetUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var user models.User
	err = Collections.Users.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Don't return the password
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Surname  string `json:"surname" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=customer admin craftsman"`
}

// CreateUser handles creating a new user
func CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Handle validation errors
		switch {
		case err.Error() == "Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag":
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		case err.Error() == "Key: 'CreateUserRequest.Role' Error:Field validation for 'Role' failed on the 'oneof' tag":
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	// Hash password first to validate password requirements
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		switch err {
		case auth.ErrPasswordTooShort:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		case auth.ErrPasswordNoUpper:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one uppercase letter"})
		case auth.ErrPasswordNoLower:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one lowercase letter"})
		case auth.ErrPasswordNoNumber:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one number"})
		case auth.ErrPasswordNoSpecial:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one special character"})
		case auth.ErrPasswordCommon:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password is too common or easily guessable"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		}
		return
	}

	// Check if email already exists
	var existingUser models.User
	err = Collections.Users.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Create new user
	user := models.User{
		Name:      req.Name,
		Surname:   req.Surname,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      models.UserRole(req.Role),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := Collections.Users.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	// Generate token pair
	tokenPair, err := auth.GenerateTokenPair(
		user.ID,
		string(user.Role),
		os.Getenv("JWT_ACCESS_SECRET"),
		os.Getenv("JWT_REFRESH_SECRET"),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// Don't return the hashed password
	user.Password = ""
	c.JSON(http.StatusCreated, gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"user": gin.H{
			"id":      user.ID.Hex(),
			"name":    user.Name,
			"surname": user.Surname,
			"email":   user.Email,
			"role":    user.Role,
		},
	})
}

// UpdateUser handles updating a user
func UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"username": user.Username,
			"email":    user.Email,
			"password": user.Password,
		},
	}

	result, err := Collections.Users.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser handles deleting a user
func DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := Collections.Users.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(router *gin.Engine) {
	users := router.Group("/users")
	{
		users.POST("", CreateUser)
		users.GET("", GetUsers)
		users.GET("/:id", GetUser)
		users.PUT("/:id", UpdateUser)
		users.DELETE("/:id", DeleteUser)
	}
}
