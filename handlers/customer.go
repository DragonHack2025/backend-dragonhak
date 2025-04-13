package handlers

import (
	"context"
	"net/http"
	"time"

	"backend-dragonhak/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var bookingCollection *mongo.Collection

// InitBookingCollection initializes the booking collection
func InitBookingCollection(mongoClient *mongo.Client) {
	bookingCollection = mongoClient.Database("dragonhak").Collection("bookings")
}

// SearchCraftsmen searches for craftsmen based on criteria
func SearchCraftsmen(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get query parameters
	specialty := c.Query("specialty")
	location := c.Query("location")
	minRating := c.Query("min_rating")

	filter := bson.M{}
	if specialty != "" {
		filter["specialties"] = specialty
	}
	if location != "" {
		filter["location"] = location
	}
	if minRating != "" {
		filter["rating"] = bson.M{"$gte": minRating}
	}

	var craftsmen []models.Craftsman
	cursor, err := Collections.Craftsmen.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var craftsman models.Craftsman
		if err := cursor.Decode(&craftsman); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		craftsmen = append(craftsmen, craftsman)
	}

	c.JSON(http.StatusOK, craftsmen)
}

// SearchWorkshops searches for available workshops
func SearchWorkshops(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get query parameters
	category := c.Query("category")
	difficulty := c.Query("difficulty")
	location := c.Query("location")
	date := c.Query("date")

	filter := bson.M{"status": "upcoming"}
	if category != "" {
		filter["category"] = category
	}
	if difficulty != "" {
		filter["difficulty"] = difficulty
	}
	if location != "" {
		filter["location"] = location
	}
	if date != "" {
		// Parse date and add to filter
		parsedDate, err := time.Parse("2006-01-02", date)
		if err == nil {
			filter["date"] = bson.M{
				"$gte": parsedDate,
				"$lt":  parsedDate.Add(24 * time.Hour),
			}
		}
	}

	var workshops []models.Workshop
	cursor, err := Collections.Workshops.Find(ctx, filter)
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

// BookWorkshop creates a booking for a workshop
func BookWorkshop(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get workshop ID from URL
	workshopID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(workshopID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workshop ID format"})
		return
	}

	// Get workshop details
	var workshop models.Workshop
	err = Collections.Workshops.FindOne(ctx, bson.M{"_id": objID}).Decode(&workshop)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workshop not found"})
		return
	}

	// Check if workshop is full
	if workshop.CurrentStudents >= workshop.MaxParticipants {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workshop is full"})
		return
	}

	// Get customer ID from the test context
	customerID := c.GetHeader("X-Customer-ID")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID is required"})
		return
	}

	customerObjID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID format"})
		return
	}

	// Create booking
	booking := models.Booking{
		ID:         primitive.NewObjectID(),
		WorkshopID: objID,
		CustomerID: customerObjID,
		Status:     models.BookingStatusConfirmed,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Insert booking
	_, err = Collections.Bookings.InsertOne(ctx, booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	// Update workshop's current students count
	_, err = Collections.Workshops.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"current_students": 1}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workshop"})
		return
	}

	// Return booking details
	c.JSON(http.StatusCreated, gin.H{
		"id":          booking.ID.Hex(),
		"workshop_id": booking.WorkshopID.Hex(),
		"customer_id": booking.CustomerID.Hex(),
		"status":      booking.Status,
		"created_at":  booking.CreatedAt,
		"updated_at":  booking.UpdatedAt,
	})
}

// GetCustomerBookings retrieves all bookings for a customer
func GetCustomerBookings(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	customerID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var bookings []models.Booking
	cursor, err := bookingCollection.Find(ctx, bson.M{"customer_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var booking models.Booking
		if err := cursor.Decode(&booking); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		bookings = append(bookings, booking)
	}

	c.JSON(http.StatusOK, bookings)
}
