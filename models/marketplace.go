package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ItemStatus represents the current status of a marketplace item
type ItemStatus string

const (
	ItemStatusActive    ItemStatus = "active"
	ItemStatusSold      ItemStatus = "sold"
	ItemStatusCancelled ItemStatus = "cancelled"
)

// ItemType represents whether the item is for sale or auction
type ItemType string

const (
	ItemTypeSale    ItemType = "sale"
	ItemTypeAuction ItemType = "auction"
)

// ItemCondition represents the condition of an item
type ItemCondition string

const (
	ItemConditionNew     ItemCondition = "new"
	ItemConditionUsed    ItemCondition = "used"
	ItemConditionAntique ItemCondition = "antique"
)

// MarketplaceItem represents an item listed in the marketplace
type MarketplaceItem struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"` // For sale items, this is the fixed price
	Type        ItemType           `json:"type" bson:"type"`
	Status      ItemStatus         `json:"status" bson:"status"`
	Images      []string           `json:"images" bson:"images"` // URLs to item images
	Category    string             `json:"category" bson:"category"`
	Condition   ItemCondition      `json:"condition" bson:"condition"`
	Location    string             `json:"location" bson:"location"`
	SellerID    primitive.ObjectID `json:"seller_id" bson:"seller_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// Auction represents an ongoing auction for an item
type Auction struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ItemID        primitive.ObjectID `json:"item_id" bson:"item_id"`
	StartingPrice float64            `json:"starting_price" bson:"starting_price"`
	CurrentPrice  float64            `json:"current_price" bson:"current_price"`
	StartTime     time.Time          `json:"start_time" bson:"start_time"`
	EndTime       time.Time          `json:"end_time" bson:"end_time"`
	Bids          []Bid              `json:"bids" bson:"bids"`
	Status        ItemStatus         `json:"status" bson:"status"`
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
