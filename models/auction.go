package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuctionItem represents an item being auctioned
type AuctionItem struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Image       string             `json:"image" bson:"image"`
	Category    string             `json:"category" bson:"category"`
	Condition   string             `json:"condition" bson:"condition"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// Auction represents an ongoing auction
type Auction struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Item          AuctionItem        `json:"item" bson:"item"`
	SellerID      primitive.ObjectID `json:"seller_id" bson:"seller_id"`
	StartingPrice float64            `json:"starting_price" bson:"starting_price"`
	CurrentPrice  float64            `json:"current_price" bson:"current_price"`
	StartTime     time.Time          `json:"start_time" bson:"start_time"`
	EndTime       time.Time          `json:"end_time" bson:"end_time"`
	LastBid       *Bid               `json:"last_bid,omitempty" bson:"last_bid,omitempty"`
	IsActive      bool               `json:"is_active" bson:"is_active"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
}

// Bid represents a bid placed in an auction
type Bid struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	AuctionID primitive.ObjectID `json:"auction_id" bson:"auction_id"`
	BidderID  primitive.ObjectID `json:"bidder_id" bson:"bidder_id"`
	Amount    float64            `json:"amount" bson:"amount"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
