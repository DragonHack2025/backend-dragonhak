package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"backend-dragonhak/models"
)

// CreateAuction handles the creation of a new auction
func CreateAuction(c *gin.Context) {
	var auction models.Auction
	if err := c.ShouldBindJSON(&auction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Validate required fields
	if auction.Item.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Auction title is required"})
		return
	}
	if auction.StartingPrice <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Starting price must be greater than 0"})
		return
	}
	if auction.EndTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be in the future"})
		return
	}

	// Set initial values
	auction.ID = primitive.NewObjectID()
	auction.Item.ID = primitive.NewObjectID()
	auction.CurrentPrice = auction.StartingPrice
	auction.IsActive = true
	auction.CreatedAt = time.Now()
	auction.UpdatedAt = time.Now()
	auction.StartTime = time.Now()

	// Get seller ID from context
	userID := c.GetString("userID")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	auction.SellerID = objID

	// Insert the auction into the database
	_, err = Collections.Auctions.InsertOne(context.Background(), auction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating auction"})
		return
	}

	c.JSON(http.StatusCreated, auction)
}

// GetAuctions retrieves all auctions with optional filters
func GetAuctions(c *gin.Context) {
	// Build filter based on query parameters
	filter := bson.M{}
	if category := c.Query("category"); category != "" {
		filter["item.category"] = category
	}
	if active := c.Query("active"); active == "true" {
		filter["is_active"] = true
		filter["end_time"] = bson.M{"$gt": time.Now()}
	}

	// Set up options for sorting and pagination
	findOptions := options.Find()
	findOptions.SetSort(bson.D{bson.E{Key: "created_at", Value: -1}})

	cursor, err := Collections.Auctions.Find(context.Background(), filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving auctions"})
		return
	}
	defer cursor.Close(context.Background())

	var auctions []models.Auction
	if err := cursor.All(context.Background(), &auctions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding auctions"})
		return
	}

	c.JSON(http.StatusOK, auctions)
}

// GetAuction retrieves a specific auction by ID
func GetAuction(c *gin.Context) {
	auctionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	var auction models.Auction
	err = Collections.Auctions.FindOne(context.Background(), bson.M{"_id": auctionID}).Decode(&auction)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Auction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving auction"})
		return
	}

	c.JSON(http.StatusOK, auction)
}

// PlaceBid handles placing a bid on an auction
func PlaceBid(c *gin.Context) {
	auctionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	var bid models.Bid
	if err := c.ShouldBindJSON(&bid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Get bidder ID from context
	userID := c.GetString("userID")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	bid.BidderID = objID
	bid.AuctionID = auctionID
	bid.ID = primitive.NewObjectID()
	bid.CreatedAt = time.Now()

	var auction models.Auction
	err = Collections.Auctions.FindOne(context.Background(), bson.M{"_id": auctionID}).Decode(&auction)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Auction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving auction"})
		return
	}

	// Validate bid
	if !auction.IsActive || time.Now().After(auction.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Auction is not active"})
		return
	}
	if bid.Amount <= auction.CurrentPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bid amount must be higher than current price"})
		return
	}
	if auction.SellerID == objID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Seller cannot bid on their own auction"})
		return
	}

	// Update auction with new bid
	update := bson.M{
		"$set": bson.M{
			"current_price": bid.Amount,
			"last_bid":      bid,
			"updated_at":    time.Now(),
		},
	}

	_, err = Collections.Auctions.UpdateOne(context.Background(), bson.M{"_id": auctionID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating auction"})
		return
	}

	// Save bid in bids collection
	_, err = Collections.Bids.InsertOne(context.Background(), bid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving bid"})
		return
	}

	c.JSON(http.StatusOK, bid)
}

// GetAuctionBids retrieves all bids for a specific auction
func GetAuctionBids(c *gin.Context) {
	auctionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{bson.E{Key: "amount", Value: -1}})

	cursor, err := Collections.Bids.Find(context.Background(), bson.M{"auction_id": auctionID}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving bids"})
		return
	}
	defer cursor.Close(context.Background())

	var bids []models.Bid
	if err := cursor.All(context.Background(), &bids); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding bids"})
		return
	}

	c.JSON(http.StatusOK, bids)
}
