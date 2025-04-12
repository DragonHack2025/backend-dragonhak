package handlers

import (
	"context"
	"net/http"
	"time"

	"backend-dragonhak/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateBadge creates a new badge
func CreateBadge(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var badge models.Badge
	if err := c.ShouldBindJSON(&badge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid badge data: " + err.Error()})
		return
	}

	// Validate required fields
	if badge.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Badge name is required"})
		return
	}

	badge.CreatedAt = time.Now()

	result, err := Collections.Badges.InsertOne(ctx, badge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create badge: " + err.Error()})
		return
	}

	badge.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, badge)
}

// AwardBadge awards a badge to a user
func AwardBadge(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID := c.Param("userId")
	badgeID := c.Param("badgeId")

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	badgeObjID, err := primitive.ObjectIDFromHex(badgeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid badge ID format"})
		return
	}

	// Get the badge
	var badge models.Badge
	err = Collections.Badges.FindOne(ctx, bson.M{"_id": badgeObjID}).Decode(&badge)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Badge not found"})
		return
	}

	// Update user's badges
	update := bson.M{
		"$addToSet": bson.M{
			"badges": badge,
		},
	}

	result, err := Collections.Users.UpdateOne(ctx, bson.M{"_id": userObjID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Badge awarded successfully"})
}

// GetUserBadges retrieves all badges for a user
func GetUserBadges(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(userID)
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

	c.JSON(http.StatusOK, user.Badges)
}

// GetBadges retrieves all available badges
func GetBadges(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var badges []models.Badge
	cursor, err := Collections.Badges.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var badge models.Badge
		if err := cursor.Decode(&badge); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		badges = append(badges, badge)
	}

	c.JSON(http.StatusOK, badges)
}
