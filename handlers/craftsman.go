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

// CreateCraftsmanProfile handles creating a new craftsman profile
func CreateCraftsmanProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var profile models.Craftsman
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if profile.Bio == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bio is required"})
		return
	}

	// Check if user exists
	var user models.User
	err := Collections.Users.FindOne(ctx, bson.M{"_id": profile.UserID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Set timestamps
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	// Insert the profile
	result, err := Collections.CraftsmanProfiles.InsertOne(ctx, profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the created profile with ID
	profile.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, profile)
}

// CreateCraft handles creating a new craft
func CreateCraft(c *gin.Context) {
	var craft models.Craft
	if err := c.ShouldBindJSON(&craft); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if craft.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	if craft.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Description is required"})
		return
	}

	if craft.Category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category is required"})
		return
	}

	if craft.Difficulty == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Difficulty is required"})
		return
	}

	if craft.Duration <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Duration must be greater than 0"})
		return
	}

	if craft.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be greater than 0"})
		return
	}

	// Set timestamps
	craft.CreatedAt = time.Now()
	craft.UpdatedAt = time.Now()

	// Insert the craft
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := Collections.Crafts.InsertOne(ctx, craft)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the created craft with ID
	craft.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, craft)
}

// CreateWorkshop creates a new workshop for a craft
func CreateWorkshop(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var workshop models.Workshop
	if err := c.ShouldBindJSON(&workshop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workshop data: " + err.Error()})
		return
	}

	// Validate required fields
	if workshop.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workshop title is required"})
		return
	}
	if workshop.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workshop description is required"})
		return
	}
	if workshop.MaxParticipants <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum number of participants must be greater than 0"})
		return
	}
	if workshop.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price cannot be negative"})
		return
	}
	if workshop.Duration <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Duration must be greater than 0"})
		return
	}

	workshop.CreatedAt = time.Now()
	workshop.UpdatedAt = time.Now()

	result, err := Collections.Workshops.InsertOne(ctx, workshop)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workshop: " + err.Error()})
		return
	}

	workshop.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, workshop)
}

// GetCraftsmanWorkshops retrieves all workshops for a craftsman
func GetCraftsmanWorkshops(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	craftsmanID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(craftsmanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var workshops []models.Workshop
	cursor, err := Collections.Workshops.Find(ctx, bson.M{"craftsman_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var workshop models.Workshop
		if err := cursor.Decode(&workshop); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		workshops = append(workshops, workshop)
	}

	c.JSON(http.StatusOK, workshops)
}

// UpdateCraftsmanProfile handles updating an existing craftsman profile
func UpdateCraftsmanProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get profile ID from URL
	profileID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Bind the update data
	var updateData models.Craftsman
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if updateData.Bio == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bio is required"})
		return
	}

	// Prepare update document
	update := bson.M{
		"$set": bson.M{
			"bio":          updateData.Bio,
			"specialties":  updateData.Specialties,
			"experience":   updateData.Experience,
			"location":     updateData.Location,
			"contact_info": updateData.ContactInfo,
			"updated_at":   time.Now(),
		},
	}

	// Update the profile
	result, err := Collections.CraftsmanProfiles.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		update,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile: " + err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No changes were made to the profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
