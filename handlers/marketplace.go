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

// CreateMarketplaceItem creates a new marketplace item
func CreateMarketplaceItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var item models.MarketplaceItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item data: " + err.Error()})
		return
	}

	// Validate required fields
	if item.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item title is required"})
		return
	}
	if item.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be greater than 0"})
		return
	}
	if item.Type != models.ItemTypeSale && item.Type != models.ItemTypeAuction {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item type"})
		return
	}

	// Set timestamps
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	item.Status = models.ItemStatusActive

	// If it's an auction item, create an auction
	if item.Type == models.ItemTypeAuction {
		auction := models.Auction{
			ItemID:        item.ID,
			StartingPrice: item.Price,
			CurrentPrice:  item.Price,
			StartTime:     time.Now(),
			EndTime:       time.Now().Add(7 * 24 * time.Hour),
			Status:        models.ItemStatusActive,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		auctionResult, err := Collections.Auctions.InsertOne(ctx, auction)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create auction: " + err.Error()})
			return
		}
		auction.ID = auctionResult.InsertedID.(primitive.ObjectID)
	}

	result, err := Collections.Marketplace.InsertOne(ctx, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item: " + err.Error()})
		return
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, item)
}

// GetMarketplaceItems retrieves all marketplace items with optional filters
func GetMarketplaceItems(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get query parameters
	category := c.Query("category")
	itemType := c.Query("type")
	status := c.Query("status")

	// Build filter
	filter := bson.M{}
	if category != "" {
		filter["category"] = category
	}
	if itemType != "" {
		filter["type"] = itemType
	}
	if status != "" {
		filter["status"] = status
	}

	var items []models.MarketplaceItem
	cursor, err := Collections.Marketplace.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item models.MarketplaceItem
		if err := cursor.Decode(&item); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, items)
}

// GetMarketplaceItem retrieves a specific marketplace item
func GetMarketplaceItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var item models.MarketplaceItem
	err = Collections.Marketplace.FindOne(ctx, bson.M{"_id": objID}).Decode(&item)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// PlaceBid places a bid on an auction item
func PlaceBid(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auctionID := c.Param("auctionId")
	auctionObjID, err := primitive.ObjectIDFromHex(auctionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID format"})
		return
	}

	var bid models.Bid
	if err := c.ShouldBindJSON(&bid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the auction
	var auction models.Auction
	err = Collections.Auctions.FindOne(ctx, bson.M{"_id": auctionObjID}).Decode(&auction)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Auction not found"})
		return
	}

	// Check if auction is still active
	if auction.Status != models.ItemStatusActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Auction is not active"})
		return
	}

	// Check if bid amount is higher than current price
	if bid.Amount <= auction.CurrentPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bid amount must be higher than current price"})
		return
	}

	// Create the bid
	bid.ID = primitive.NewObjectID()
	bid.AuctionID = auctionObjID
	bid.CreatedAt = time.Now()

	// Update auction with new bid
	update := bson.M{
		"$set": bson.M{
			"current_price": bid.Amount,
			"updated_at":    time.Now(),
		},
		"$push": bson.M{
			"bids": bid,
		},
	}

	_, err = Collections.Auctions.UpdateOne(ctx, bson.M{"_id": auctionObjID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place bid"})
		return
	}

	// Save the bid
	_, err = Collections.Bids.InsertOne(ctx, bid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save bid"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bid placed successfully",
		"bid":     bid,
	})
}

// GetAuctionBids retrieves all bids for an auction
func GetAuctionBids(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auctionID := c.Param("auctionId")
	auctionObjID, err := primitive.ObjectIDFromHex(auctionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID format"})
		return
	}

	var bids []models.Bid
	cursor, err := Collections.Bids.Find(ctx, bson.M{"auction_id": auctionObjID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var bid models.Bid
		if err := cursor.Decode(&bid); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		bids = append(bids, bid)
	}

	c.JSON(http.StatusOK, bids)
}
