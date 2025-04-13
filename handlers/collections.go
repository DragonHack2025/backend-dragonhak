package handlers

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection interface defines the methods needed for MongoDB operations
type Collection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
}

// Collections holds all MongoDB collections
var Collections struct {
	Users             Collection
	CraftsmanProfiles Collection
	Crafts            Collection
	Workshops         Collection
	Badges            Collection
	Auctions          Collection
	Bids              Collection
	Bookings          Collection
}

// InitCollections initializes all collections
func InitCollections(db *mongo.Database) {
	Collections.Users = db.Collection("users")
	Collections.CraftsmanProfiles = db.Collection("craftsman_profiles")
	Collections.Crafts = db.Collection("crafts")
	Collections.Workshops = db.Collection("workshops")
	Collections.Badges = db.Collection("badges")
	Collections.Auctions = db.Collection("auctions")
	Collections.Bids = db.Collection("bids")
	Collections.Bookings = db.Collection("bookings")
}
