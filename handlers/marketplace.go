package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"backend-dragonhak/models"
	"backend-dragonhak/utils"
)

// CreateAuction handles the creation of a new auction
func CreateAuction(w http.ResponseWriter, r *http.Request) {
	var auction models.Auction
	if err := json.NewDecoder(r.Body).Decode(&auction); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate required fields
	if auction.Item.Title == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Auction title is required")
		return
	}
	if auction.StartingPrice <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Starting price must be greater than 0")
		return
	}
	if auction.EndTime.Before(time.Now()) {
		utils.RespondWithError(w, http.StatusBadRequest, "End time must be in the future")
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
	userID := r.Context().Value("userID").(primitive.ObjectID)
	auction.SellerID = userID

	// Insert the auction into the database
	_, err := Collections.Auctions.InsertOne(context.Background(), auction)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating auction")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, auction)
}

// GetAuctions retrieves all auctions with optional filters
func GetAuctions(w http.ResponseWriter, r *http.Request) {
	collection := Collections.Auctions

	// Build filter based on query parameters
	filter := bson.M{}
	if category := r.URL.Query().Get("category"); category != "" {
		filter["item.category"] = category
	}
	if active := r.URL.Query().Get("active"); active == "true" {
		filter["is_active"] = true
		filter["end_time"] = bson.M{"$gt": time.Now()}
	}

	// Set up options for sorting and pagination
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving auctions")
		return
	}
	defer cursor.Close(context.Background())

	var auctions []models.Auction
	if err := cursor.All(context.Background(), &auctions); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error decoding auctions")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, auctions)
}

// GetAuction retrieves a specific auction by ID
func GetAuction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	auctionID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid auction ID")
		return
	}

	collection := Collections.Auctions
	var auction models.Auction
	err = collection.FindOne(context.Background(), bson.M{"_id": auctionID}).Decode(&auction)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.RespondWithError(w, http.StatusNotFound, "Auction not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving auction")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, auction)
}

// PlaceBid handles placing a bid on an auction
func PlaceBid(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	auctionID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid auction ID")
		return
	}

	var bid models.Bid
	if err := json.NewDecoder(r.Body).Decode(&bid); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Get bidder ID from context
	userID := r.Context().Value("userID").(primitive.ObjectID)
	bid.BidderID = userID
	bid.AuctionID = auctionID
	bid.ID = primitive.NewObjectID()
	bid.CreatedAt = time.Now()

	collection := Collections.Auctions
	var auction models.Auction
	err = collection.FindOne(context.Background(), bson.M{"_id": auctionID}).Decode(&auction)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.RespondWithError(w, http.StatusNotFound, "Auction not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving auction")
		return
	}

	// Validate bid
	if !auction.IsActive || time.Now().After(auction.EndTime) {
		utils.RespondWithError(w, http.StatusBadRequest, "Auction is not active")
		return
	}
	if bid.Amount <= auction.CurrentPrice {
		utils.RespondWithError(w, http.StatusBadRequest, "Bid amount must be higher than current price")
		return
	}
	if auction.SellerID == userID {
		utils.RespondWithError(w, http.StatusBadRequest, "Seller cannot bid on their own auction")
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

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": auctionID}, update)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating auction")
		return
	}

	// Save bid in bids collection
	bidsCollection := Collections.Bids
	_, err = bidsCollection.InsertOne(context.Background(), bid)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error saving bid")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, bid)
}

// GetAuctionBids retrieves all bids for a specific auction
func GetAuctionBids(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	auctionID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid auction ID")
		return
	}

	collection := Collections.Bids
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "amount", Value: -1}})

	cursor, err := collection.Find(context.Background(), bson.M{"auction_id": auctionID}, findOptions)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving bids")
		return
	}
	defer cursor.Close(context.Background())

	var bids []models.Bid
	if err := cursor.All(context.Background(), &bids); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error decoding bids")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, bids)
}
