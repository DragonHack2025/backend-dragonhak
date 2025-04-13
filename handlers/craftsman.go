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
	"golang.org/x/crypto/bcrypt"
)

// CreateCraftsmanProfile creates a new craftsman profile
func CreateCraftsmanProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var request struct {
		Username    string                    `json:"username" binding:"required"`
		Email       string                    `json:"email" binding:"required,email"`
		Password    string                    `json:"password" binding:"required,min=8"`
		Bio         string                    `json:"bio"`
		Experience  int                       `json:"experience"`
		Rating      float64                   `json:"rating" binding:"required,min=0,max=5"`
		Location    string                    `json:"location" binding:"required"`
		ContactInfo models.ContactInformation `json:"contact_info"`
		IsVerified  bool                      `json:"is_verified"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existingUser models.User
	err := Collections.Users.FindOne(ctx, bson.M{"email": request.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// Create user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username:  request.Username,
		Email:     request.Email,
		Password:  string(hashedPassword),
		Role:      "craftsman",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userResult, err := Collections.Users.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	userID := userResult.InsertedID.(primitive.ObjectID)

	// Create craftsman profile
	craftsman := models.Craftsman{
		UserID:      userID,
		Bio:         request.Bio,
		Experience:  request.Experience,
		Rating:      request.Rating,
		Location:    request.Location,
		ContactInfo: request.ContactInfo,
		IsVerified:  request.IsVerified,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	craftsmanResult, err := Collections.Craftsmen.InsertOne(ctx, craftsman)
	if err != nil {
		// If craftsman creation fails, delete the user
		Collections.Users.DeleteOne(ctx, bson.M{"_id": userID})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create craftsman profile"})
		return
	}

	// Generate token pair
	tokenPair, err := auth.GenerateTokenPair(
		userID,
		string(user.Role),
		os.Getenv("JWT_ACCESS_SECRET"),
		os.Getenv("JWT_REFRESH_SECRET"),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"user": gin.H{
			"id":    userID.Hex(),
			"email": user.Email,
			"role":  user.Role,
		},
		"craftsman": gin.H{
			"id":           craftsmanResult.InsertedID.(primitive.ObjectID).Hex(),
			"bio":          craftsman.Bio,
			"experience":   craftsman.Experience,
			"rating":       craftsman.Rating,
			"location":     craftsman.Location,
			"contact_info": craftsman.ContactInfo,
			"is_verified":  craftsman.IsVerified,
		},
	})
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
	result, err := Collections.Craftsmen.UpdateOne(
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

// GetCraftsmen retrieves all craftsmen
func GetCraftsmen(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var craftsmen []models.Craftsman
	cursor, err := Collections.Craftsmen.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &craftsmen); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, craftsmen)
}

// GetCraftsman retrieves a specific craftsman by ID
func GetCraftsman(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var craftsman models.Craftsman
	err = Collections.Craftsmen.FindOne(ctx, bson.M{"_id": objectID}).Decode(&craftsman)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Craftsman not found"})
		return
	}

	c.JSON(http.StatusOK, craftsman)
}

// UpdateCraftsman updates a craftsman's profile
func UpdateCraftsman(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var craftsman models.Craftsman
	if err := c.ShouldBindJSON(&craftsman); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"bio":          craftsman.Bio,
			"experience":   craftsman.Experience,
			"rating":       craftsman.Rating,
			"location":     craftsman.Location,
			"contact_info": craftsman.ContactInfo,
			"is_verified":  craftsman.IsVerified,
			"updated_at":   time.Now(),
		},
	}

	result, err := Collections.Craftsmen.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Craftsman not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Craftsman updated successfully"})
}

// DeleteCraftsman deletes a craftsman's profile
func DeleteCraftsman(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := Collections.Craftsmen.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Craftsman not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Craftsman deleted successfully"})
}
